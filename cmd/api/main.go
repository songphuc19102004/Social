package main

import (
	"log"

	env "github.com/songphuc19102004/social/internal"
	"github.com/songphuc19102004/social/internal/db"
	"github.com/songphuc19102004/social/internal/store"
)

const version string = "0.0.1"

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:123123@localhost:5431/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("db connected")

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store:  *store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
