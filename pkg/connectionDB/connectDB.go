package connDB

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

func NewConnetPostgres() (*sqlx.DB, error) {
	cfg := ConfigParamDb{}
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	if err := envconfig.Process("db", &cfg); err != nil {
		return nil, err
	}
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBname, cfg.Sslmode))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
