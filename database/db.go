package database

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/ximofam/chess-online/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Driver       string
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
}

const (
	defaultMaxOpenConns = 25
	defaultMaxIdleConns = 10
)

func (c *Config) applyDefaults() {
	if c.MaxOpenConns <= 0 {
		c.MaxOpenConns = defaultMaxOpenConns
	}

	if c.MaxIdleConns <= 0 {
		c.MaxIdleConns = defaultMaxIdleConns
	}
}

func New(cfg Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch cfg.Driver {
	case "mysql":
		db, err = newMySQL(&cfg)
	default:
		return nil, errors.New("unsupported database driver")
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	cfg.applyDefaults()

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Printf("connect to %s successfully", cfg.Driver)

	return db, nil
}

func newMySQL(cfg *Config) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // output
		logger.Config{
			SlowThreshold:             time.Second, // log query chậm
			LogLevel:                  logger.Info, // mức log
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	return gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: newLogger,
	})
}

// GormDatabase is a wrapper for the gorm framework.
type GormDatabase struct {
	DB *gorm.DB
}

// Close closes the gorm database connection.
func (d *GormDatabase) Close() {
	sqldb, err := d.DB.DB()
	if err != nil {
		return
	}
	sqldb.Close()
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
	)
}
