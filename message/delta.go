package message

import (
	"fmt"
	"strings"
)

// DeltaType is the type of mutation for the sql to take
type DeltaType uint8

const (
	// Insert is for an INSERT SQL statement
	Insert = iota + 1
	// Update is for an UPDATE SQL statement
	Update
	// Delete is for a DELETE SQL statement
	Delete
)

// InsertDelta is the delta type for inserts
type InsertDelta struct {
	Record
	Table string
}

// UpdateDelta is the delta type for updates
type UpdateDelta struct {
	IDRecord        // IDRecord as we need to identify this record in the update statement
	Changes  Record // for holding old values
	Table    string
}

// DeleteDelta is the delta type for updates
type DeleteDelta struct {
	IDRecord // IDRecord as we need to identify this record in the delete statement
	Table    string
}

// NewInsertDelta is an insert change record related to a table
func NewInsertDelta(r Record, table string) *InsertDelta {
	return &InsertDelta{
		Record: r,
		Table:  table,
	}
}

// NewUpdateDelta is an insert change record related to a table
func NewUpdateDelta(r IDRecord, table string) *UpdateDelta {
	return &UpdateDelta{
		IDRecord: r,
		Changes:  NewRecord(),
		Table:    table,
	}
}

// NewDeleteDelta is a delete change record related to a table.
func NewDeleteDelta(r IDRecord, table string) *DeleteDelta {
	return &DeleteDelta{
		IDRecord: r,
		Table:    table,
	}
}

// GetSQL implements the SQLGetter interface
func (d InsertDelta) GetSQL() string {
	vals := ""
	sep := ""
	for range d.GetKeys() {
		vals += sep + "?"
		sep = ","
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, d.Table, strings.Join(d.GetKeys(), ","), vals)
}

// GetSQL implements the SQLGetter interface
func (d UpdateDelta) GetSQL() string {
	cols := strings.Join(d.GetKeys(), "=? AND ") + "=?"
	id := strings.Join(d.GetIDKeys(), "=? AND ") + "=?"
	return fmt.Sprintf(`UPDATE %s SET %s WHERE %s`, d.Table, cols, id)
}

// GetSQL implements the SQLGetter interface
func (d DeleteDelta) GetSQL() string {
	id := strings.Join(d.GetIDKeys(), "=? AND ") + "=?"
	return fmt.Sprintf(`DELETE FROM %s WHERE %s`, d.Table, id)
}

// GetArgs implements the ArgGetter interface
func (d InsertDelta) GetArgs() []interface{} {
	return d.GetVals()
}

// GetArgs implements the ArgGetter interface
func (d UpdateDelta) GetArgs() []interface{} {
	return append(d.GetVals(), d.GetIDVals())
}

// GetArgs implements the ArgGetter interface
func (d DeleteDelta) GetArgs() []interface{} {
	return d.GetIDVals()
}
