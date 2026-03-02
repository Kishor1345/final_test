// Package credentials provides database connection pooling utilities.
//
// This package is responsible for initializing, exposing, and
// gracefully closing the PostgreSQL connection pool used across
// the HRMODULE application.
//
// The pool is initialized once during application startup and
// reused by all APIs to ensure efficient and safe database access.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:29-01-2026
//
// Last Modified By:
//
// Last Modified Date:
package credentials

import (
	"database/sql"
	"fmt"
	"sync"

	// pgx is used as the PostgreSQL driver and integrates with database/sql
	// while providing better performance and compatibility.
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	// db holds the singleton PostgreSQL connection pool.
	db *sql.DB

	// once ensures the pool is initialized only once.
	once sync.Once

	// dbErr stores any error encountered during pool initialization.
	dbErr error
)

// InitPostgresPool initializes the PostgreSQL connection pool using
// the provided connection string.
//
// This function must be called exactly once during application startup
// (typically from main.go) before any database operations are performed.
//
// The pool uses Go's default database/sql settings without any
// custom tuning or time-based limits.
func InitPostgresPool(connStr string) error {

	once.Do(func() {
		db, dbErr = sql.Open("pgx", connStr)
		if dbErr != nil {
			dbErr = fmt.Errorf("failed to open postgres DB: %w", dbErr)
			return
		}

		// Validate database connectivity
		if err := db.Ping(); err != nil {
			dbErr = fmt.Errorf("postgres DB ping failed: %w", err)
		}
	})

	return dbErr
}

// GetDB returns the initialized PostgreSQL *sql.DB connection pool.
//
// This function panics if the pool has not been initialized,
// which helps catch configuration errors early during startup.
func GetDB() *sql.DB {
	if db == nil {
		panic("PostgreSQL DB pool not initialized. Call InitPostgresPool() in main.go")
	}
	return db
}

// ClosePostgresPool gracefully closes the PostgreSQL connection pool.
//
// This should be called during application shutdown to ensure all
// open connections are properly released.
func ClosePostgresPool() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
