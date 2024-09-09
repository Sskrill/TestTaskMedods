package connDB

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestNewConnetPostgres_Mock(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error create sqlmock DB: %v", err)
	}
	defer mockDB.Close()
	mock.ExpectPing()
	db := sqlx.NewDb(mockDB, "postgres")
	if err := db.Ping(); err != nil {
		t.Errorf("Error ping: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error except: %v", err)
	}
}
