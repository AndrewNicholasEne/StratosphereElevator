package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/AndrewNicholasEne/StratosphereElevator/internal/db"
	"github.com/AndrewNicholasEne/StratosphereElevator/internal/graph"
	"github.com/AndrewNicholasEne/StratosphereElevator/internal/graph/generated"
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

	gqlSrv := handler.New(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: &graph.Resolver{StacksService: stackSvc},
		}),
	)

	gqlSrv.AddTransport(transport.GET{})
	gqlSrv.AddTransport(transport.POST{})

	gqlSrv.SetParserTokenLimit(1_000_000)

	mux := http.NewServeMux()
	h.Register(mux)

	mux.Handle("/graphql", gqlSrv) // accepts POST (and GET for queries)
	mux.Handle("/graphiql", playground.Handler("GraphQL", "/graphql"))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
