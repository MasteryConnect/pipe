package message

import (
	"fmt"
	"strings"
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
	idkeys := d.GetIDKeys()
	id := strings.Join(idkeys, "=? AND ") + "=?"
	// get the cols without the id
	var writeCols []string

	for _, col := range d.GetKeys() {
		exclude := false
		for _, idkey := range idkeys {
			if col == idkey {
				exclude = true
				break // skip this col as it exists in the id
			}
		}
		if !exclude {
			writeCols = append(writeCols, col)
		}
	}
	cols := strings.Join(writeCols, "=?, ") + "=?"

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
// and returns the args for the query excluding the
// args related to the ID keys.
func (d UpdateDelta) GetArgs() []interface{} {
	// get the non-ID keys
	vals := []interface{}{}
	for _, key := range GetNonIDKeys(d) {
		if val, ok := d.Get(key); ok {
			vals = append(vals, val)
		} else {
			return nil // bail because we didn't find a value
		}
	}
	return vals
}

// GetArgs implements the ArgGetter interface
func (d DeleteDelta) GetArgs() []interface{} {
	return d.GetIDVals()
}
