package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// we are not using dbConfig struct as parameter because internal should not know outside its scope
func New(addr string, maxOpenConns int, maxIdleConns int, maxIdleTimeout string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTimeout)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
