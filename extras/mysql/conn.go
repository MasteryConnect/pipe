package mysql

import (
	"github.com/jmoiron/sqlx"
	// include the mysql sql driver
	_ "github.com/go-sql-driver/mysql"
)

// Conn is a wrapper around a mysql connection to handle
// creating the connection and shutting it down.
// This allows you to just set a DSN and this will handle connecting to it.
type Conn struct {
	Dsn       string // connection string (required)
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
		m.DB, err = sqlx.Open("mysql", m.Dsn)
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
