package sql

import (
	"github.com/jmoiron/sqlx"
	// include the sql driver of your choice in your own code to use the driver here
	// check this out for a list of options : https://github.com/golang/go/wiki/SQLDrivers
	// Ex:
	//_ "github.com/go-sql-driver/mysql"
	//_ "github.com/lib/pq"
	//_ "github.com/mattn/go-sqlite3"
)

// Conn is a wrapper around a connection to handle
// creating the connection and shutting it down.
// This allows you to just set a DSN and this will handle connecting to it.
type Conn struct {
	DSN       string // connection string (required)
	Driver    string // used to make a new connection if DB is nil
	DB        *sqlx.DB
	customDB  bool
	connected bool // prevent multiple Open calls
}

// Open will connect to the database if we are handed a DSN instead of an *sqlx.DB .
func (m *Conn) Open() error {
	if m.connected {
		return nil
	}

	var err error
	if m.DB == nil {
		m.DB, err = sqlx.Open(m.Driver, m.DSN)
		if err != nil {
			return err
		}
	} else {
		m.customDB = true
	}

	m.connected = true
	return m.DB.Ping()
}

// Close will close the connection to the database if we created the connection in Startup.
func (m *Conn) Close() error {
	if m.customDB {
		return nil
	}
	return m.DB.Close()
}
