package setup

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	goccyyaml "github.com/goccy/go-yaml"
	"gopkg.in/yaml.v3"

	"github.com/crowdsecurity/crowdsec/pkg/csconfig"
	"github.com/crowdsecurity/crowdsec/pkg/cwhub"
)

// AcquisDocument is created from a SetupItem. It represents a single YAML document, and can be part of a multi-document file.
type AcquisDocument struct {
	AcquisFilename string
	DataSource     map[string]interface{}
}

func decodeSetup(input []byte, fancyErrors bool) (Setup, error) {
	ret := Setup{}

	// parse with goccy to have better error messages in many cases
	dec := goccyyaml.NewDecoder(bytes.NewBuffer(input), goccyyaml.Strict())

	if err := dec.Decode(&ret); err != nil {
		if fancyErrors {
			return ret, fmt.Errorf("%v", goccyyaml.FormatError(err, true, true))
		}
		// XXX errors here are multiline, should we just print them to stderr instead of logging?
		return ret, fmt.Errorf("%v", err)
	}

	// parse again because goccy is not strict enough anyway
	dec2 := yaml.NewDecoder(bytes.NewBuffer(input))
	dec2.KnownFields(true)

	if err := dec2.Decode(&ret); err != nil {
		return ret, fmt.Errorf("while unmarshaling setup file: %w", err)
	}

	return ret, nil
}

// InstallHubItems installs the objects recommended in a setup file.
func InstallHubItems(csConfig *csconfig.Config, input []byte, dryRun bool) error {
	setupEnvelope, err := decodeSetup(input, false)
	if err != nil {
		return err
	}

	if err := csConfig.LoadHub(); err != nil {
		return fmt.Errorf("loading hub: %w", err)
	}

	if err := cwhub.SetHubBranch(); err != nil {
		return fmt.Errorf("setting hub branch: %w", err)
	}

	if err := cwhub.GetHubIdx(csConfig.Hub); err != nil {
		return fmt.Errorf("getting hub index: %w", err)
	}

	for _, setupItem := range setupEnvelope.Setup {
		forceAction := false
		downloadOnly := false
		install := setupItem.Install

		if install == nil {
			continue
		}

		if len(install.Collections) > 0 {
			for _, collection := range setupItem.Install.Collections {
				if dryRun {
					fmt.Println("dry-run: would install collection", collection)

					continue
				}

				if err := cwhub.InstallItem(csConfig, collection, cwhub.COLLECTIONS, forceAction, downloadOnly); err != nil {
					return fmt.Errorf("while installing collection %s: %w", collection, err)
				}
			}
		}

		if len(install.Parsers) > 0 {
			for _, parser := range setupItem.Install.Parsers {
				if dryRun {
					fmt.Println("dry-run: would install parser", parser)

					continue
				}

				if err := cwhub.InstallItem(csConfig, parser, cwhub.PARSERS, forceAction, downloadOnly); err != nil {
					return fmt.Errorf("while installing parser %s: %w", parser, err)
				}
			}
		}

		if len(install.Scenarios) > 0 {
			for _, scenario := range setupItem.Install.Scenarios {
				if dryRun {
					fmt.Println("dry-run: would install scenario", scenario)

					continue
				}

				if err := cwhub.InstallItem(csConfig, scenario, cwhub.SCENARIOS, forceAction, downloadOnly); err != nil {
					return fmt.Errorf("while installing scenario %s: %w", scenario, err)
				}
			}
		}

		if len(install.PostOverflows) > 0 {
			for _, postoverflow := range setupItem.Install.PostOverflows {
				if dryRun {
					fmt.Println("dry-run: would install postoverflow", postoverflow)

					continue
				}

				if err := cwhub.InstallItem(csConfig, postoverflow, cwhub.PARSERS_OVFLW, forceAction, downloadOnly); err != nil {
					return fmt.Errorf("while installing postoverflow %s: %w", postoverflow, err)
				}
			}
		}
	}

	return nil
}

// marshalAcquisDocuments creates the monolithic file, or itemized files (if a directory is provided) with the acquisition documents.
func marshalAcquisDocuments(ads []AcquisDocument, toDir string) (string, error) {
	var sb strings.Builder

	dashTerminator := false

	disclaimer := `
#
# This file was automatically generated by "cscli setup datasources".
# You can modify it by hand, but will be responsible for its maintenance.
# To add datasources or logfiles, you can instead write a new configuration
# in the directory defined by acquisition_dir.
#

`

	if toDir == "" {
		sb.WriteString(disclaimer)
	} else {
		_, err := os.Stat(toDir)
		if os.IsNotExist(err) {
			return "", fmt.Errorf("directory %s does not exist", toDir)
		}
	}

	for _, ad := range ads {
		out, err := goccyyaml.MarshalWithOptions(ad.DataSource, goccyyaml.IndentSequence(true))
		if err != nil {
			return "", fmt.Errorf("while encoding datasource: %w", err)
		}

		if toDir != "" {
			if ad.AcquisFilename == "" {
				return "", fmt.Errorf("empty acquis filename")
			}

			fname := filepath.Join(toDir, ad.AcquisFilename)
			fmt.Println("creating", fname)

			f, err := os.Create(fname)
			if err != nil {
				return "", fmt.Errorf("creating acquisition file: %w", err)
			}
			defer f.Close()

			_, err = f.WriteString(disclaimer)
			if err != nil {
				return "", fmt.Errorf("while writing to %s: %w", ad.AcquisFilename, err)
			}

			_, err = f.Write(out)
			if err != nil {
				return "", fmt.Errorf("while writing to %s: %w", ad.AcquisFilename, err)
			}

			f.Sync()

			continue
		}

		if dashTerminator {
			sb.WriteString("---\n")
		}

		sb.Write(out)

		dashTerminator = true
	}

	return sb.String(), nil
}

// Validate checks the validity of a setup file.
func Validate(input []byte) error {
	_, err := decodeSetup(input, true)
	if err != nil {
		return err
	}

	return nil
}

// DataSources generates the acquisition documents from a setup file.
func DataSources(input []byte, toDir string) (string, error) {
	setupEnvelope, err := decodeSetup(input, false)
	if err != nil {
		return "", err
	}

	ads := make([]AcquisDocument, 0)

	filename := func(basename string, ext string) string {
		if basename == "" {
			return basename
		}

		return basename + ext
	}

	for _, setupItem := range setupEnvelope.Setup {
		datasource := setupItem.DataSource

		basename := ""
		if toDir != "" {
			basename = "setup." + setupItem.DetectedService
		}

		if datasource == nil {
			continue
		}

		ad := AcquisDocument{
			AcquisFilename: filename(basename, ".yaml"),
			DataSource:     datasource,
		}
		ads = append(ads, ad)
	}

	return marshalAcquisDocuments(ads, toDir)
}
