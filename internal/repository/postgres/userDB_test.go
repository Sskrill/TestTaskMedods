package pgSql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sskrill/TestTaskMedods/internal/domain"
	"github.com/jmoiron/sqlx"
)

func TestUserDB_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error create sqlmock db: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	userDB := NewUserDB(sqlxDB)

	mock.ExpectExec("INSERT INTO users").
		WithArgs("MedodsTest", "MedodsTest@mail.ru", "123456789", time.Now()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	ctx := context.Background()
	user := domain.User{
		Name:         "MedodsTest",
		Email:        "MedodsTest@mail.ru",
		Password:     "123456789",
		RegisteredAt: time.Now(),
	}
	err = userDB.Create(ctx, user)
	if err != nil {
		t.Errorf("Expect err=nil, error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations succes: %v", err)
	}
}

func TestUserDB_GetByParams(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error create sqlmock db: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	userDB := NewUserDB(sqlxDB)

	rows := sqlmock.NewRows([]string{"name", "email", "registered_at", "id"}).
		AddRow("MedodsTest", "MedodsTest@mail.ru", time.Now(), 1)
	mock.ExpectQuery("SELECT name,email,registered_at,id FROM users WHERE email=\\$1 AND password=\\$2").
		WithArgs("MedodsTest@mail.ru", "123456789").
		WillReturnRows(rows)

	ctx := context.Background()
	user, err := userDB.GetByParams(ctx, "MedodsTest@mail.ru", "123456789")
	if err != nil {
		t.Errorf("Expect err=nil, error: %v", err)
	}

	if user.Name != "MedodsTest" || user.Email != "MedodsTest@mail.ru" {
		t.Errorf("Expect Name 'MedodsTest' && email 'MedodsTest@mail.ru', results %s && %s", user.Name, user.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations succes: %v", err)
	}
}

func TestUserDB_GetEmailById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error create sqlmock db: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	userDB := NewUserDB(sqlxDB)

	rows := sqlmock.NewRows([]string{"email"}).AddRow("MedodsTest@mail.ru")
	mock.ExpectQuery("SELECT email FROM users WHERE id=\\$1").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	ctx := context.Background()
	email, err := userDB.GetEmailById(ctx, 1)
	if err != nil {
		t.Errorf("Expect err=nil, error: %v", err)
	}

	if email != "MedodsTest@mail.ru" {
		t.Errorf("Expect email 'MedodsTest@mail.ru', result %s", email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations succes: %v", err)
	}
}
