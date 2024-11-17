package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewConnection(host, port, user, password, dbname string) (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
        user, password, host, port, dbname)
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %v", err)
    }

    db.SetMaxIdleConns(10)
    db.SetMaxOpenConns(100)
    db.SetConnMaxLifetime(time.Hour)

    err = db.Ping()
    if err != nil {
        return nil, fmt.Errorf("error connecting to the database: %v", err)
    }

    return db, nil
}