package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"

	db "todo-app/internal/db/sqlc"
	"todo-app/internal/handlers"
)

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errMissingEnv("DATABASE_URL")
	}

	dbConn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return err
	}

	if err := dbConn.Ping(); err != nil {
		return err
	}

	queries := db.New(dbConn)

	r := chi.NewRouter()
	handlers.RegisterTodoRoutes(r, queries)

	log.Println("Server running on :8080")
	return http.ListenAndServe(":8080", r)
}

type errMissingEnv string

func (e errMissingEnv) Error() string {
	return "missing environment variable: " + string(e)
}
