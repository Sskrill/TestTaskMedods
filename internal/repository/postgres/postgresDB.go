package pgSql

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type PostgresDB interface {
	Migration()
}

type postgresDB struct {
	db *sqlx.DB
}

func NewPostgresDB(db *sqlx.DB) *postgresDB { return &postgresDB{db} }

func (pDB *postgresDB) Migration() {
	pDB.db.MustExec(`CREATE TABLE users(
  id SERIAL PRIMARY KEY,
name VARCHAR(25) NOT NULL,
password VARCHAR(250) NOT NULL,
email VARCHAR(50) NOT NULL,
registered_at TIMESTAMP NOT NULL
);

CREATE TABLE tokens(
id SERIAL PRIMARY KEY,
user_id INT NOT NULL,
r_token VARCHAR(250) NOT NULL,
a_token VARCHAR(250) NOT NULL,
uuid VARCHAR(250) NOT NULL,
expires_at TIMESTAMP NOT NULL,
ip_address VARCHAR(255) NOT NULL,
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE);
`)

}

func (pDB *postgresDB) CloseDB() {
	err := pDB.db.Close()
	if err != nil {
		log.Println(err)
	}

	log.Println("Connection to Postgres closed.")
}
