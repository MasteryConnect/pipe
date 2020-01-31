package sql

import (
	"strconv"
	"time"

	// include the mysql sql driver
	"github.com/MasteryConnect/pipe/message"
	"github.com/jmoiron/sqlx"
)

const (
	// MysqlDatetimeFormat is the Mysql date time format.
	MysqlDatetimeFormat = "2006-01-02 15:04:05"
)

// Querier gets records from a query or table
type Querier Conn

// T will take in records and use them in a sql query.
func (m Querier) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	c := Conn(m)
	err := c.Open()
	if err != nil {
		errs <- err
	}
	defer func() {
		err := c.Close()
		if err != nil {
			errs <- err
		}
	}()

	db := Querier(c)

	for msg := range in {
		rows, err := db.Run(msg)
		if err != nil {
			errs <- err
		} else {
			extractRows(rows.(*sqlx.Rows), out, errs)
		}
	}
}

// Run actually does the query again the database and
// implements the InlineFunc interface.
func (m Querier) Run(msg interface{}) (interface{}, error) {
	switch v := msg.(type) {
	case message.SQLGetter:
		if a, ok := v.(message.ArgsGetter); ok && a.GetArgs() != nil {
			return m.DB.Queryx(v.GetSQL(), a.GetArgs()...)
		}
		return m.DB.Queryx(v.GetSQL())
	default:
		return m.DB.Queryx(message.String(v))
	}
}

// UnpackRows is a transformer (Tfunc) that iterates over the sqlx.Rows and
// sends downstream each row as an OrderedRecord.
func UnpackRows(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for m := range in {
		extractRows(m.(*sqlx.Rows), out, errs)
	}
}

func extractRows(rows *sqlx.Rows, out chan<- interface{}, errs chan<- error) (cnt int) {
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
				row[k] = castByColumnType(string(b), colTypes[k])
			}
		}

		out <- message.NewRecordFromMSI(row).SetKeyOrder(cols...)
	}
	return
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