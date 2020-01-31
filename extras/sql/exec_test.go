package sql_test

import (
	"github.com/MasteryConnect/pipe/extras/sql"
	"github.com/MasteryConnect/pipe/line"
	"github.com/MasteryConnect/pipe/message"
)

func ExampleExec() {
	line.New().SetP(func(out chan<- interface{}, errs chan<- error) {

		// using a string directly as the query
		out <- "INSERT INTO foo (name) VALUES ('bar')"

		// using message.Query
		out <- message.Query{SQL: "INSERT INTO foo (name) VALUES ('bar')"}

		// using message.Query with args (using sqlx under the hood for arg matching)
		query := message.Query{SQL: "INSERT INTO foo (name) VALUES (?)"}
		query.Args = append(query.Args, "bar")
		out <- query

		// using message.Delta
		record := message.NewRecordFromMSI(map[string]interface{}{"name": "bar"})
		delta := message.InsertDelta{Table: "foo", Record: record}
		out <- delta

	}).Add(
		sql.Exec{Driver: "postgres", DSN: "postgres://user:pass@localhost/dbname?sslmode=verify-full"}.T,
	).Run()
}
