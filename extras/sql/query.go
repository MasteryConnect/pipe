package sql

import (
	"strconv"
	"time"

	"github.com/MasteryConnect/pipe/message"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	// MysqlDatetimeFormat is the Mysql date time format.
	MysqlDatetimeFormat = "2006-01-02 15:04:05"
)

// Query runs an SQL query on a db connection
type Query Conn

// T will take in records and use them in a sql query.
func (m Query) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	c := Conn(m)
	if err := c.Open(); err != nil {
		errs <- err
	}
	defer func() {
		err := c.Close()
		if err != nil {
			errs <- err
		}
	}()

	db := Query(c)

	for msg := range in {
		results, err := db.I(msg)
		if err != nil {
			errs <- err
		} else {
			rows := results.(*sqlx.Rows)
			ExtractRecordsFromRows(rows, func(row message.OrderedRecord, err error) error {
				if err != nil {
					errs <- err
				} else {
					out <- row
				}
				return nil // only return an error here if we want the iterating to stop
			})
			rows.Close()
		}
	}
}

// I actually does the query again the database and
// implements the InlineTfunc interface.
func (m Query) I(msg interface{}) (interface{}, error) {
	switch v := msg.(type) {
	case message.ToSQLer:
		sql, args := v.ToSQL()
		return m.DB.Queryx(sql, args...)
	default:
		return m.DB.Queryx(message.String(v))
	}
}

// ErrStopExtract is the error to be used by a rowHandler to indicate that
// the iterating over the rows results should not continue.
var ErrStopExtract = errors.New("stop iterating over the *sqlx.Rows")

// ExtractRecordsFromRows iterates over the *sqlx.Rows and calls the handler for each
// or error if something went wrong. The row handler can return it's own error
// in which case the iterating still stop. If the error is not the ErrStopExtract
// error, then it will be pass on through as the error of the overall func call.
// Keep in mind you will need to make sure the rows passed in are close elsewhere.
// (Ex: rows.Close() ) This func will not close the rows.
func ExtractRecordsFromRows(rows *sqlx.Rows, rowHandler func(message.OrderedRecord, error) error) error {
	colTypes := make(map[string]string)

	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	types, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	// fill a map of the types by column name
	for i, t := range types {
		colTypes[cols[i]] = t.DatabaseTypeName()
	}

	var row map[string]interface{}
	for rows.Next() {
		row = make(map[string]interface{})

		err = rows.MapScan(row)
		if err != nil {
			err = rowHandler(nil, err)
			if err != nil {
				if err == ErrStopExtract {
					break // stop but don't pass the "Stop" error up
				} else {
					return err // stop and pass the error on up
				}
			}
			continue
		}

		for k, v := range row {
			if i, ok := v.(int64); ok {
				row[k] = i
			} else if t, ok := v.(time.Time); ok {
				row[k] = t
			} else if b, ok := v.([]byte); ok {
				row[k] = castByColumnType(string(b), colTypes[k])
			}
		}

		err = rowHandler(message.NewRecordFromMSI(row).SetKeyOrder(cols...), nil)
		if err != nil {
			if err == ErrStopExtract {
				break // stop but don't pass the "Stop" error up
			} else {
				return err // stop and pass the error on up
			}
		}
	}
	return nil
}

func castByColumnType(input, columnType string) (output interface{}) {
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
