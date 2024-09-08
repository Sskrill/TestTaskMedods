package app

import (
	"fmt"
	"github.com/Sskrill/TestTaskMedods/internal/api"
	pgSql "github.com/Sskrill/TestTaskMedods/internal/repository/postgres"
	"github.com/Sskrill/TestTaskMedods/internal/service"
	connDB "github.com/Sskrill/TestTaskMedods/pkg/connectionDB"
	"github.com/Sskrill/TestTaskMedods/pkg/hasher"
	"log"
	"net/http"
	"os"
)

func Run() {
	pg, err := connDB.NewConnetPostgres()
	if err != nil {
		log.Fatal(err)
	}
	//pgDB := pgSql.NewPostgresDB(pg) // убрать коментарии при первом запуске
	//pgDB.Migration()
	userDB := pgSql.NewUserDB(pg)
	tokenDB := pgSql.NewTokensDB(pg)
	passwordHasher := hasher.NewHasher(os.Getenv("Salt"))

	srvc := service.NewUserService(passwordHasher, tokenDB, userDB, []byte(os.Getenv("Secret")))
	handler := api.NewHandler(srvc)
	server := &http.Server{Addr: os.Getenv("Port"), Handler: handler.CreateRouter()}
	fmt.Println("Server started || Сервер зарущен")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
