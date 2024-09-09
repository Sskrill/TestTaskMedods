package pgSql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sskrill/TestTaskMedods/internal/domain"
	"github.com/jmoiron/sqlx"
)

func TestTokensDB_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error create sqlmock db: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	tokensDB := NewTokensDB(sqlxDB)
	rows := sqlmock.NewRows([]string{"id", "user_id", "r_token", "expires_at", "ip_address", "a_token", "uuid"}).
		AddRow(1, 10, "some_r_token", time.Now(), "127.0.0.1", "some_a_token", "uuid123")
	mock.ExpectQuery("SELECT id,user_id,r_token,expires_at,ip_address,a_token,uuid FROM tokens WHERE r_token = \\$1").
		WithArgs("some_r_token").
		WillReturnRows(rows)
	ctx := context.Background()
	token, err := tokensDB.Get(ctx, "some_r_token")
	if err != nil {
		t.Errorf("Expect err=nil,error: %v", err)
	}
	if token.RToken != "some_r_token" {
		t.Errorf("Expect r_token = some_r_token ,result = %s", token.RToken)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations succes: %v", err)
	}
}

func TestTokensDB_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error create sqlmock db: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	tokensDB := NewTokensDB(sqlxDB)
	mock.ExpectExec("INSERT INTO tokens").
		WithArgs(10, "some_r_token", time.Now(), "127.0.0.1", "some_a_token", "uuid123").
		WillReturnResult(sqlmock.NewResult(1, 1))
	ctx := context.Background()
	token := domain.Tokens{
		UserId:    10,
		RToken:    "some_r_token",
		ExpiresAt: time.Now(),
		IpAddr:    "127.0.0.1",
		AToken:    "some_a_token",
		Uuid:      "uuid123",
	}
	err = tokensDB.Create(ctx, token)
	if err != nil {
		t.Errorf("Expect err=nil,error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations succes: %v", err)
	}
}
func TestTokensDB_DeleteToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error create sqlmock db: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	tokensDB := NewTokensDB(sqlxDB)
	mock.ExpectExec("DELETE FROM tokens WHERE user_id=\\$1").
		WithArgs(int64(10)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	ctx := context.Background()
	err = tokensDB.DeleteToken(ctx, 10)
	if err != nil {
		t.Errorf("Expect err=nil,error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations succes: %v", err)
	}
}

func TestTokensDB_GetCoupleOfTokens(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error create sqlmock db: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	tokensDB := NewTokensDB(sqlxDB)
	rows := sqlmock.NewRows([]string{"a_token", "r_token"}).
		AddRow("access_token_value", "refresh_token_value")
	mock.ExpectQuery("SELECT a_token,r_token FROM tokens WHERE uuid=\\$1").
		WithArgs("uuid123").
		WillReturnRows(rows)
	ctx := context.Background()
	aT, rT, err := tokensDB.GetCoupleOfTokens(ctx, "uuid123")
	if err != nil {
		t.Errorf("Expect err=nil,error: %v", err)
	}
	if aT != "access_token_value" || rT != "refresh_token_value" {
		t.Errorf("Expect aT = access_token_value && rT = refresh_token_value, result %s && %s", aT, rT)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectations succes: %v", err)
	}
}
