package shared

import "database/sql"

// DBTX is an interface that allows us to pass both *sql.DB and *sql.Tx
// to our repository methods, enabling transactional support.
type DBTX interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
