package config

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"testing"
)

func TestDatabase(t *testing.T) {
	// Save original env vars
	originalUser := os.Getenv("MYSQL_USER")
	originalPass := os.Getenv("MYSQL_PASSWORD")
	originalHost := os.Getenv("MYSQL_HOST")
	originalPort := os.Getenv("MYSQL_PORT")

	// Set test env vars
	os.Setenv("MYSQL_USER", "test_user")
	os.Setenv("MYSQL_PASSWORD", "test_password")
	os.Setenv("MYSQL_HOST", "localhost")
	os.Setenv("MYSQL_PORT", "3306")

	// Cleanup function to restore original env vars
	defer func() {
		os.Setenv("MYSQL_USER", originalUser)
		os.Setenv("MYSQL_PASSWORD", originalPass)
		os.Setenv("MYSQL_HOST", originalHost)
		os.Setenv("MYSQL_PORT", originalPort)
	}()

	t.Run("Database Connection Success", func(t *testing.T) {
		db := Database()
		if db == nil {
			t.Error("Expected database connection, got nil")
		}
		defer db.Close()

		err := db.Ping()
		if err != nil {
			t.Errorf("Database connection failed: %v", err)
		}
	})
}

func TestDatabaseMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Test create database
	mock.ExpectExec("CREATE DATABASE gotodo").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("USE gotodo").WillReturnResult(sqlmock.NewResult(0, 0))
	
	// Test create table
	mock.ExpectExec("CREATE TABLE todos").WillReturnResult(sqlmock.NewResult(0, 0))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestDatabaseEnvVars(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "All env vars present",
			envVars: map[string]string{
				"MYSQL_USER":     "test_user",
				"MYSQL_PASSWORD": "test_pass",
				"MYSQL_HOST":     "localhost",
				"MYSQL_PORT":     "3306",
			},
			wantErr: false,
		},
		{
			name: "Missing user",
			envVars: map[string]string{
				"MYSQL_PASSWORD": "test_pass",
				"MYSQL_HOST":     "localhost",
				"MYSQL_PORT":     "3306",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear existing env vars
			os.Clearenv()

			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			if !tt.wantErr {
				db := Database()
				if db == nil {
					t.Error("Expected database connection, got nil")
				}
			}
		})
	}
}

func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	db := Database()
	
	cleanup := func() {
		_, err := db.Exec("DROP DATABASE IF EXISTS gotodo")
		if err != nil {
			t.Errorf("Error cleaning up test database: %v", err)
		}
		db.Close()
	}

	return db, cleanup
}

func TestCleanup(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	// Verify database exists
	var dbName string
	err := db.QueryRow("SELECT DATABASE()").Scan(&dbName)
	if err != nil {
		t.Fatalf("Failed to get database name: %v", err)
	}

	if dbName != "gotodo" {
		t.Errorf("Expected database name 'gotodo', got %s", dbName)
	}
}