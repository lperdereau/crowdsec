// Code generated by entc, DO NOT EDIT.

package machine

import (
	"time"
)

const (
	// Label holds the string label denoting the machine type in the database.
	Label = "machine"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldLastPush holds the string denoting the last_push field in the database.
	FieldLastPush = "last_push"
	// FieldMachineId holds the string denoting the machineid field in the database.
	FieldMachineId = "machine_id"
	// FieldPassword holds the string denoting the password field in the database.
	FieldPassword = "password"
	// FieldIpAddress holds the string denoting the ipaddress field in the database.
	FieldIpAddress = "ip_address"
	// FieldScenarios holds the string denoting the scenarios field in the database.
	FieldScenarios = "scenarios"
	// FieldVersion holds the string denoting the version field in the database.
	FieldVersion = "version"
	// FieldIsValidated holds the string denoting the isvalidated field in the database.
	FieldIsValidated = "is_validated"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// EdgeAlerts holds the string denoting the alerts edge name in mutations.
	EdgeAlerts = "alerts"
	// Table holds the table name of the machine in the database.
	Table = "machines"
	// AlertsTable is the table that holds the alerts relation/edge.
	AlertsTable = "alerts"
	// AlertsInverseTable is the table name for the Alert entity.
	// It exists in this package in order to avoid circular dependency with the "alert" package.
	AlertsInverseTable = "alerts"
	// AlertsColumn is the table column denoting the alerts relation/edge.
	AlertsColumn = "machine_alerts"
)

// Columns holds all SQL columns for machine fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldLastPush,
	FieldMachineId,
	FieldPassword,
	FieldIpAddress,
	FieldScenarios,
	FieldVersion,
	FieldIsValidated,
	FieldStatus,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// UpdateDefaultCreatedAt holds the default value on update for the "created_at" field.
	UpdateDefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultLastPush holds the default value on creation for the "last_push" field.
	DefaultLastPush func() time.Time
	// UpdateDefaultLastPush holds the default value on update for the "last_push" field.
	UpdateDefaultLastPush func() time.Time
	// ScenariosValidator is a validator for the "scenarios" field. It is called by the builders before save.
	ScenariosValidator func(string) error
	// DefaultIsValidated holds the default value on creation for the "isValidated" field.
	DefaultIsValidated bool
)
