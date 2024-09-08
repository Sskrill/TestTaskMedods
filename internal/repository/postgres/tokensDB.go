package pgSql

import (
	"context"
	"github.com/Sskrill/TestTaskMedods/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TokensDB struct {
	db *sqlx.DB
}

func NewTokensDB(db *sqlx.DB) *TokensDB { return &TokensDB{db: db} }

func (tDB *TokensDB) Get(ctx context.Context, rToken string) (domain.Tokens, error) {
	var token domain.Tokens
	err := tDB.db.QueryRow("SELECT id,user_id,r_token,expires_at,ip_address,a_token,uuid FROM tokens WHERE r_token = $1", rToken).
		Scan(&token.Id, &token.UserId, &token.RToken, &token.ExpiresAt, &token.IpAddr, &token.AToken, &token.Uuid)
	return token, err
}
func (tDB *TokensDB) Create(ctx context.Context, token domain.Tokens) error {
	_, err := tDB.db.Exec("INSERT INTO tokens (user_id,r_token,expires_at,ip_address,a_token,uuid) VALUES($1,$2,$3,$4,$5,$6)",
		token.UserId, token.RToken, token.ExpiresAt, token.IpAddr, token.AToken, token.Uuid)
	return err
}

func (tDB *TokensDB) DeleteToken(ctx context.Context, userId int64) error {
	_, err := tDB.db.Exec("DELETE FROM tokens WHERE user_id=$1", userId)
	return err
}

func (tDB *TokensDB) GetCoupleOfTokens(ctx context.Context, guid string) (aT, rT string, err error) {
	err = tDB.db.QueryRow("SELECT a_token,r_token FROM tokens WHERE uuid=$1", guid).Scan(&aT, &rT)
	return aT, rT, err
}
