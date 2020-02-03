package sql

import (
	dbsql "database/sql"
	"fmt"

	"github.com/MasteryConnect/pipe/message"
	"github.com/pkg/errors"
)

// ErrInvalidSQLQueryType is the error for when the query passed in is
// isn't one of the know types.
var ErrInvalidSQLQueryType = errors.New("invalid SQL query type")

// Exec executes a query
type Exec Conn

// T will take in records and use them in a sql query.
// The input message is expected to implement the fmt.Stringer interface
// and string is expected to be the SQL to run.
func (m Exec) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
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

	db := Exec(c) // use c to call I() on to keep the db connection

	for msg := range in {
		result, err := db.I(msg)

		if err != nil {
			errs <- err
		}
		if result != nil {
			out <- message.SQLResult{Result: result.(dbsql.Result)}
		}
	}
}

// I implements the InlineFunc interface. It is context aware
// if the incoming message is a mysql.Query and has the context set.
func (m Exec) I(msg interface{}) (interface{}, error) {
	switch v := msg.(type) {
	case message.ToSQLer:
		sql, args := v.ToSQL()
		if c, cok := msg.(message.ContextGetter); cok {
			ctx := c.GetContext()
			if ctx != nil {
				return m.DB.ExecContext(ctx, sql, args...)
			}
		}
		return m.DB.Exec(sql, args...)
	case string:
		return m.DB.Exec(v)
	case fmt.Stringer:
		return m.DB.Exec(v.String())
	}
	return nil, errors.Wrapf(ErrInvalidSQLQueryType, "got unknown type %T", msg)
}
