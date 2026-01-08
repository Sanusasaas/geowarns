package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct{
	DB *gorm.DB
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConnection(config *Config) (*Repository, error){
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName,config.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err!=nil {
		return nil, err
	}
	return &Repository{DB: db}, nil
}

func (r *Repository) Close() error {
	sqlDB, err := r.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}