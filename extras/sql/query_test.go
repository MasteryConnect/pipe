package sql_test

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/MasteryConnect/pipe/extras/sql"
	"github.com/MasteryConnect/pipe/line"
	"github.com/MasteryConnect/pipe/message"

	// in memory db driver for testing
	"github.com/jmoiron/sqlx"
	_ "github.com/proullon/ramsql/driver"
)

func ExampleQuery() {
	line.New().SetP(func(out chan<- interface{}, errs chan<- error) {

		// using a string directly as the query
		out <- "SELECT * FROM foo"

		// using message.Query
		out <- message.Query{SQL: "SELECT * FROM foo"}

		// using message.Query with args (using sqlx under the hood for arg matching)
		query := message.Query{SQL: "SELECT * FROM foo WHERE name=? AND email=?"}
		query.Args = append(query.Args, "bar", "bar@foo.com")
		out <- query

	}).Add(
		sql.Query{Driver: "postgres", DSN: "postgres://user:pass@localhost/dbname?sslmode=verify-full"}.T,
	).Run()
}

func TestQuery(t *testing.T) {
	setup := []string{
		"CREATE TABLE users (id BIGSERIAL PRIMARY KEY, name TEXT, email TEXT, age INT);",
		"INSERT INTO users (name,email,age) VALUES ('alice','alice@foo.com',30);",
		"INSERT INTO users (name,email,age) VALUES ('bob','bob@foo.com',31);",
		"INSERT INTO users (name,email,age) VALUES ('charlie','charlie@foo.com',32);",
	}

	db, err := sqlx.Open("ramsql", "TestQuery")
	if err != nil {
		t.Fatalf("ramsql.Open: Error: %s\n", err)
		return
	}
	defer db.Close()

	for _, step := range setup {
		if _, err = db.Exec(step); err != nil {
			t.Fatalf("ramsql.Exec: Error: %s\n", err)
		}
	}

	// use Query to run a query
	queryTranformer := sql.Query{DB: db}
	resultsI, err := queryTranformer.I("SELECT age,name,email FROM users")
	if err != nil {
		t.Fatalf("sql.Query: Error: %s\n", err)
	}

	// Examine the results
	var names, emails []string
	var ages []int
	results := resultsI.(*sqlx.Rows)
	defer results.Close()

	sql.ExtractRecordsFromRows(results, func(row message.OrderedRecord, err error) error {
		// check the column order
		wantOrder := []string{"age", "name", "email"}
		if !reflect.DeepEqual(wantOrder, row.GetKeys()) {
			t.Errorf("want %v got %v", wantOrder, row.GetKeys())
		}

		name, _ := row.Get("name")
		email, _ := row.Get("email")
		ageStr, _ := row.Get("age")
		age, _ := strconv.Atoi(fmt.Sprint(ageStr))
		names = append(names, name.(string))
		emails = append(emails, email.(string))
		ages = append(ages, age)
		return nil
	})

	if len(names) != 3 {
		t.Errorf("want %d got %d", 3, len(names))
		return
	}
	want := []string{"alice", "bob", "charlie"}
	if !reflect.DeepEqual(want, names) {
		t.Errorf("want %v got %v", want, names)
	}
	want = []string{"alice@foo.com", "bob@foo.com", "charlie@foo.com"}
	if !reflect.DeepEqual(want, emails) {
		t.Errorf("want %v got %v", want, emails)
	}
	wantAges := []int{30, 31, 32}
	if !reflect.DeepEqual(wantAges, ages) {
		t.Errorf("want %v got %v", wantAges, ages)
	}
}
