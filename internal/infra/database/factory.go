package database

import (
	"os"
	"youmeet/internal/adapters/repositories"
)

func NewDBClient() (repositories.DBClient, error) {
	env := os.Getenv("DB_TYPE")
	
	switch env {
	case "sqlite":
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "youmeet.db"
		}
		return NewSQLiteClient(dbPath)
		
	case "postgres":
		fallthrough
	default:
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			dsn = "host=localhost user=youmeet password=youmeet dbname=youmeet port=5432 sslmode=disable"
		}
		return NewPostgresClient(dsn)
	}
}