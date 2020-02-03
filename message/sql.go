package message

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// ToSQLer defines the SQL func to extract the SQL from a message.
type ToSQLer interface {
	ToSQL() (string, []interface{})
}

// Query is an SQL query that has the arguments in a sparate slice.
// This allows for the DB.Exec() call to properly handle the args.
type Query struct {
	SQL     string
	Args    []interface{}
	Context context.Context
}

// SQLResult wraps the sql.Result and adds a String func for convenient output of the results.
type SQLResult struct {
	sql.Result
}

// String implements the fmt.Stringer interface
func (r SQLResult) String() string {
	rowCnt, err := r.RowsAffected()
	return fmt.Sprintf("%d rows affected (err: %v)", rowCnt, err)
}

// String converts the query to a completed SQL string ready to run.
// This generally should only be used in testing and debugging as
// the mysql.Exec will properly handle the args separatly from the SQL.
func (q Query) String() string {
	escaped := make([]interface{}, len(q.Args))
	for i, v := range q.Args {
		escaped[i] = strings.Replace(fmt.Sprintf("%v", v), "'", "''", -1)
	}
	return fmt.Sprintf(strings.Replace(q.SQL, "?", "'%v'", -1), escaped...)
}

// GetContext implements the message.ContextGetter interface
func (q Query) GetContext() context.Context {
	return q.Context
}

// ToSQL implements the message.SQLGetter interface
func (q Query) ToSQL() (string, []interface{}) {
	return q.SQL, q.Args
}
