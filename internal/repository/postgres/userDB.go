package pgSql

import (
	"context"
	"github.com/Sskrill/TestTaskMedods/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserDB struct {
	db *sqlx.DB
}

func NewUserDB(db *sqlx.DB) *UserDB { return &UserDB{db: db} }

func (uDB *UserDB) Create(ctx context.Context, user domain.User) error {
	_, err := uDB.db.Exec("INSERT INTO users (name,email,password,registered_at) VALUES ($1,$2,$3,$4)",
		user.Name, user.Email, user.Password, user.RegisteredAt)
	return err
}
func (uDB *UserDB) GetByParams(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	err := uDB.db.QueryRow("SELECT name,email,registered_at,id FROM users WHERE email=$1 AND password=$2", email, password).
		Scan(&user.Name, &user.Email, &user.RegisteredAt, &user.Id)
	return user, err
}
func (uDB *UserDB) GetEmailById(ctx context.Context, id int64) (string, error) {
	var email string
	err := uDB.db.QueryRow("SELECT email FROM users WHERE id=$1", id).
		Scan(&email)
	return email, err
}
