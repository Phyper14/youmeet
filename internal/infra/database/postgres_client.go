package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"youmeet/internal/adapters/repositories"
)

type PostgresClient struct {
	db *gorm.DB
}

func NewPostgresClient(dsn string) (repositories.DBClient, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	
	return &PostgresClient{db: db}, nil
}

func (p *PostgresClient) Create(value interface{}) error {
	return p.db.Create(value).Error
}

func (p *PostgresClient) First(dest interface{}, conds ...interface{}) error {
	return p.db.First(dest, conds...).Error
}

func (p *PostgresClient) Find(dest interface{}, conds ...interface{}) error {
	return p.db.Find(dest, conds...).Error
}

func (p *PostgresClient) Where(query interface{}, args ...interface{}) repositories.DBClient {
	return &PostgresClient{db: p.db.Where(query, args...)}
}

func (p *PostgresClient) Delete(value interface{}, conds ...interface{}) error {
	return p.db.Delete(value, conds...).Error
}

func (p *PostgresClient) AutoMigrate(dst ...interface{}) error {
	return p.db.AutoMigrate(dst...)
}