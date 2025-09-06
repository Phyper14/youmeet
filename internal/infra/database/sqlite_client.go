package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"youmeet/internal/adapters/repositories"
)

type SQLiteClient struct {
	db *gorm.DB
}

func NewSQLiteClient(dbPath string) (repositories.DBClient, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	
	return &SQLiteClient{db: db}, nil
}

func (s *SQLiteClient) Create(value interface{}) error {
	return s.db.Create(value).Error
}

func (s *SQLiteClient) First(dest interface{}, conds ...interface{}) error {
	return s.db.First(dest, conds...).Error
}

func (s *SQLiteClient) Find(dest interface{}, conds ...interface{}) error {
	return s.db.Find(dest, conds...).Error
}

func (s *SQLiteClient) Where(query interface{}, args ...interface{}) repositories.DBClient {
	return &SQLiteClient{db: s.db.Where(query, args...)}
}

func (s *SQLiteClient) Delete(value interface{}, conds ...interface{}) error {
	return s.db.Delete(value, conds...).Error
}

func (s *SQLiteClient) AutoMigrate(dst ...interface{}) error {
	return s.db.AutoMigrate(dst...)
}