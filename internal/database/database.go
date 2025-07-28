package database

import (
	"EverythingSuckz/fsb/internal/types"
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes the SQLite database
func InitDatabase(log *zap.Logger) error {
	log = log.Named("database")
	defer log.Sugar().Info("Initialized database")

	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Database file path
	dbPath := filepath.Join(dataDir, "fsb_stats.db")

	// Configure GORM logger
	gormLogger := logger.New(
		&GormLogWriter{log: log},
		logger.Config{
			LogLevel: logger.Info,
		},
	)

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate tables
	err = db.AutoMigrate(&types.Stats{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// GormLogWriter implements io.Writer for GORM logging
type GormLogWriter struct {
	log *zap.Logger
}

func (w *GormLogWriter) Write(p []byte) (n int, err error) {
	w.log.Debug(string(p))
	return len(p), nil
} 