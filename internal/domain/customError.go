package domain

import "errors"

var (
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

type CustomErrorResponse struct {
	Message string `json:"ErrorMessage"`
}
