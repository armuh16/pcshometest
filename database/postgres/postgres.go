package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"pcstakehometest/config"
	"pcstakehometest/package/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Gorm *gorm.DB
	Sql  *sql.DB
}

func NewPostgres(log *logger.LogRus) *DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Get().Postgres.Host,
		config.Get().Postgres.Port,
		config.Get().Postgres.Username,
		config.Get().Postgres.Password,
		config.Get().Postgres.DBName,
		config.Get().Postgres.SSLMode)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("NewPostgres", err.Error())
	}

	sqldb, err := gormDB.DB()
	if err != nil {
		log.Fatalf("NewPostgres", err.Error())
	}

	if err := sqldb.Ping(); err != nil {
		log.Fatalf("NewPostgres", err.Error())
	}

	sqldb.SetMaxOpenConns(100)
	sqldb.SetMaxIdleConns(10)
	sqldb.SetConnMaxIdleTime(300 * time.Second)
	sqldb.SetConnMaxLifetime(time.Duration(300 * time.Second))
	return &DB{
		Sql:  sqldb,
		Gorm: gormDB,
	}
}
