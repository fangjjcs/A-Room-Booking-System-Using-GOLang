package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB : holds the database connection pool
type DB struct{
	SQL *sql.DB
}

var dbConn = &DB{}


const maxOpenDbConn = 10
const maxIdelDbConn = 5
const maxDbLiveConn = 5* time.Minute

// create database pool
func ConnectSQL(dsn string) (*DB, error){
	d, err := NewDatabase(dsn)
	if err != nil{
		panic(err)
	}
	// Limitation of db connection (pool)
	d.SetConnMaxIdleTime(maxOpenDbConn)
	d.SetConnMaxIdleTime(maxIdelDbConn)
	d.SetConnMaxLifetime(maxDbLiveConn)

	dbConn.SQL = d
	testDB(d)
	if err != nil{
		return nil, err
	}
	return dbConn, nil
}

// try to ping database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err!=nil{
		return err
	}
	return nil
}

// create a new database
func NewDatabase(dsn string) (*sql.DB, error){
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err!=nil{
		return nil, err
	}

	return db, nil
}