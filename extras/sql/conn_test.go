package sql_test

import (
	"github.com/MasteryConnect/pipe/extras/sql"
)

// ExampleConn shows how to create a database connection
// that is used by the Query and Exec transformers.
// Conn is rarely used directly on it's own, but it wrapped
// by Exec and Query and used by them.
func ExampleConn() {
	c := sql.Conn{Driver: "postgres", DSN: "postgres://user:pass@localhost/dbname?sslmode=verify-full"}
	err := c.Open()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// now the connection is ready to use
}
