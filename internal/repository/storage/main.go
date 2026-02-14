package storage

import (
	"fmt"
	"graph-interview/internal/cfg"
	"graph-interview/internal/domain"
	"graph-interview/pkg/logger"
	"time"

	postgresDrv "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type DB struct {
	DB *gorm.DB
}

func NewDB(cfg *cfg.DatabaseConfig) (*DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormLogger.Default.LogMode(cfg.GormLogLevel),
	}

	gormDb, err := gorm.Open(postgresDrv.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database :%w", err)
	}
	// Set the pool size
	sqlDB, err := gormDb.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying *sql.DB: %w", err)
	}
	if cfg.MaxOpenConnectionCount == 0 {
		cfg.MaxOpenConnectionCount = 15
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnectionCount) // Maximum number of open connections
	sqlDB.SetMaxIdleConns(5)                          // Maximum number of idle connections
	sqlDB.SetConnMaxIdleTime(time.Second * 15)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	if sqlDB.Ping() != nil {
		return nil, fmt.Errorf("connection to the database couldn't be established!")
	}

	res := &DB{
		DB: gormDb,
	}

	if err := res.migration(); err != nil {
		return nil, fmt.Errorf("auto migration failed: %w", err)

	}
	return res, nil
}

func (p *DB) migration() error {
	err := p.DB.AutoMigrate(
		&domain.User{},
		&domain.UserSession{},
		&domain.Task{},
	)
	if err == nil {
		logger.Logger.Info("database migration successfully done")
	}
	return err
}
