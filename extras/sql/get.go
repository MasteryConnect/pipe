package sql

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"
	"time"

	// include the mysql sql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Get gets records from a query or table
type Get struct {
	Conn
	SQL   string // (required)
	Table string

	PageSize int
	OrderBy  string
	BodyCol  string // the column to put as the body of the message (blank is all as json)
}

// P starts sourcing the data for the pipeline from a mysql table
func (m Get) P(out chan<- interface{}, errs chan<- error) {
	err := m.Open()
	if err != nil {
		errs <- err
		return
	}
	defer func() {
		err := m.Close()
		if err != nil {
			errs <- err
		}
	}()
	m.runQuery(m.SQL, nil, out, errs)
}

// T will take in records and use them in a sql query.
// The input can be a single record or a batch of records. The predefined
// template functions can be used to extract individual keys from the metadata
// record, or the message body as a single string value. Or if this is a batch
// the keys can be used to extract one columns data as a comma seperated list
// e.g. for use in an IN clause.
//
// The output by default will send one resulting record in one out message.
// Otherwise this can be configured to send all resulting records from one
func (m Get) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	err := m.Open()
	if err != nil {
		errs <- err
	}
	defer func() {
		err := m.Close()
		if err != nil {
			errs <- err
		}
	}()

	sqlTmpl, err := template.New("mysql_get-transform").Parse(m.SQL)
	if err != nil {
		errs <- fmt.Errorf("Error parsing the sql template: %v", err)
	}

	for inMsg := range in {
		var buf bytes.Buffer
		err := sqlTmpl.Execute(&buf, &inMsg)
		if err != nil {
			errs <- fmt.Errorf("Error creating sql from template and values: %v", err)
		}
		m.runQuery(buf.String(), &inMsg, out, errs)
	}
}

func (m *Get) castByColumnType(input, columnType string) (output interface{}) {
	switch columnType {
	case "INT", "MEDIUMINT":
		ival, err := strconv.ParseInt(input, 10, 32)
		if err == nil {
			output = int32(ival)
		}
	case "TINYINT":
		ival, err := strconv.ParseInt(input, 10, 8)
		if err == nil {
			output = int8(ival)
		}
	case "SMALLINT":
		ival, err := strconv.ParseInt(input, 10, 16)
		if err == nil {
			output = int16(ival)
		}
	case "BIGINT":
		ival, err := strconv.ParseInt(input, 10, 64)
		if err == nil {
			output = ival
		}
	case "DATETIME":
		t, err := time.Parse(MysqlDatetimeFormat, input)
		if err == nil {
			output = t
		} else {
			output = input
		}
	default:
		output = input
	}
	return
}

func (m *Get) runQuery(sqlQuery string, inMsg interface{}, out chan<- interface{}, errs chan<- error) {
	rows := m.query(0, sqlQuery, errs)
	lastCount, lastID := m.process(rows, out, errs)
	rows.Close()

	for m.PageSize > 0 && lastCount > 0 {
		rows := m.query(lastID, sqlQuery, errs)
		if rows != nil {
			lastCount, lastID = m.process(rows, out, errs)
			rows.Close()
		}
	}
}

func (m *Get) query(id interface{}, sqlQuery string, errs chan<- error) *sqlx.Rows {
	var err error
	var rows *sqlx.Rows
	if sqlQuery != "" {
		rows, err = m.DB.Queryx(sqlQuery)
	} else if m.Table != "" {
		if m.PageSize == 0 {
			rows, err = m.DB.Queryx(fmt.Sprintf("SELECT * FROM %s", m.Table))
		} else {
			if m.OrderBy == "" {
				m.OrderBy = "id"
			}
			rows, err = m.DB.Queryx(fmt.Sprintf("SELECT * FROM %s WHERE %s > %d ORDER BY %s ASC LIMIT %d", m.Table, m.OrderBy, id, m.OrderBy, m.PageSize))
		}
	}
	if err != nil {
		errs <- err
	}

	return rows
}

func (m *Get) process(rows *sqlx.Rows, out chan<- interface{}, errs chan<- error) (cnt int, lastID interface{}) {
	cnt = 0
	colTypes := make(map[string]string)

	cols, err := rows.Columns()
	if err != nil {
		errs <- err
	}
	types, err := rows.ColumnTypes()
	if err != nil {
		errs <- err
	}

	// fill a map of the types by column name
	for i, t := range types {
		colTypes[cols[i]] = t.DatabaseTypeName()
	}

	var row map[string]interface{}
	for rows.Next() {
		row = make(map[string]interface{})
		cnt++

		err = rows.MapScan(row)
		if err != nil {
			errs <- err
		}

		for k, v := range row {
			if i, ok := v.(int64); ok {
				row[k] = i
			} else if t, ok := v.(time.Time); ok {
				row[k] = t
			} else if b, ok := v.([]byte); ok {
				row[k] = m.castByColumnType(string(b), colTypes[k])
			}
		}

		out <- row
		lastID = row[m.OrderBy]
	}
	return
}
