package app

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sychd/banking-auth/domain"
	"github.com/sychd/banking-auth/logger"
	"github.com/sychd/banking-auth/service"
	"log"
	"net/http"
	"os"
	"time"
)

func createDbClient(dbUrl string) *sqlx.DB {
	client, err := sqlx.Open("mysql", dbUrl)
	if err != nil {
		panic(err)
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client
}
func Start() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("godotenv error: %s", err)
		return
	}

	router := mux.NewRouter()
	dbClient := createDbClient(os.Getenv("CLEARDB_DATABASE_URL"))
	authRepository := domain.NewAuthRepository(dbClient)
	ah := AuthHandler{service.NewLoginService(authRepository, domain.GetRolePermissions())}

	router.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)
	router.HandleFunc("/auth/register", ah.NotImplementedHandler).Methods(http.MethodPost)
	router.HandleFunc("/auth/refresh", ah.Refresh).Methods(http.MethodPost)
	router.HandleFunc("/auth/verify", ah.Verify).Methods(http.MethodGet)

	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	logger.Info(fmt.Sprintf("Starting OAuth server on %s:%s ...", address, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))
}
