package api

import (
	"context"
	"github.com/Sskrill/TestTaskMedods/internal/domain"
	"github.com/gorilla/mux"
	"net/http"
)

type Service interface {
	ParseToken(ctx context.Context, token string) (int64, error)
	RefreshTokens(ctx context.Context, refreshToken, ipAddr string) (string, string, error)
	SignIn(ctx context.Context, inp domain.SignInInput, ipAddr string) (string, string, error)
	SignUp(ctx context.Context, param domain.SignUpInput) error
	GetTokensByGUID(ctx context.Context, guid string) (aT, rT string, err error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler { return &Handler{service: service} }

func (h *Handler) CreateRouter() *mux.Router {
	router := mux.NewRouter()
	auth := router.PathPrefix("/auth").Subrouter()
	{

		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodGet)
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
		auth.HandleFunc("/refresh", h.refreshTokens).Methods(http.MethodGet)
	}
	medods := router.PathPrefix("/medods").Subrouter()
	{
		medods.Use(h.authMiddleware)
		medods.HandleFunc("/guid/{guid}", h.getTokens).Methods(http.MethodGet)

	}
	return router
}
