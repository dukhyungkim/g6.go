package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

func IsSupportedEngines(engine string) bool {
	switch strings.ToLower(engine) {
	case "sqlite", "postgresql", "mysql":
		return true
	}
	return false
}

var instance *sql.DB

func NewDB(engine string) error {
	if instance != nil {
		return nil
	}

	db, err := sql.Open(engine, "sqlite3.db")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	instance = db
	return nil
}
