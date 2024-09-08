package domain

import "time"

type Tokens struct {
	Id        int64
	UserId    int64
	AToken    string
	RToken    string
	Uuid      string
	ExpiresAt time.Time
	IpAddr    string
}
type CooupleOfTokens struct {
	AToken string `json:"access_token"`
	RToken string `json:"refresh_token"`
}
