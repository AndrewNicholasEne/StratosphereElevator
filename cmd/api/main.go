package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/AndrewNicholasEne/StratosphereElevator/internal/db"
	httpapi "github.com/AndrewNicholasEne/StratosphereElevator/internal/http"
	"github.com/AndrewNicholasEne/StratosphereElevator/internal/stacks"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dbUrl := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries := db.New(pool)
	stackSvc := stacks.New(queries, logger)
	h := httpapi.NewStacksHTTP(stackSvc)

	mux := http.NewServeMux()
	h.Register(mux)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
