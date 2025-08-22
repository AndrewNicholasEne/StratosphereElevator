package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/AndrewNicholasEne/StratosphereElevator/internal/db"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	must(err)
	defer pool.Close()
	must(pool.Ping(ctx))

	q := db.New(pool)

	// 1) Create
	slug := fmt.Sprintf("demo-%d", time.Now().UnixNano())
	created, err := q.CreateStack(ctx, db.CreateStackParams{
		ID:   uuid.New(),
		Name: "Demo Stack",
		Slug: slug,
	})
	must(err)
	fmt.Printf("created: %s  %s  %s\n", created.ID, created.Slug, created.CreatedAt.UTC().Format(time.RFC3339))

	// 2) List (non-archived)
	list, err := q.ListStacks(ctx, db.ListStacksParams{
		Column1: false,
		Column2: 10,
		Column3: 0,
	})
	must(err)
	fmt.Println("list (non-archived):")
	for _, s := range list {
		if s.ArchivedAt.Valid {
			fmt.Printf("  %s  %-16s  archived %s\n", s.ID, s.Slug, s.ArchivedAt.Time.UTC().Format(time.RFC3339))
		} else {
			fmt.Printf("  %s  %-16s  created %s\n", s.ID, s.Slug, s.CreatedAt.UTC().Format(time.RFC3339))
		}
	}

	// 3) Get by slug
	got, err := q.GetStackBySlug(ctx, slug)
	must(err)
	fmt.Printf("getBySlug: %s  %s\n", got.ID, got.Slug)

	// 4) Archive
	archived, err := q.ArchiveStack(ctx, created.ID)
	must(err)
	fmt.Printf("archived: %s  at=%s  valid=%v\n",
		archived.ID, archived.ArchivedAt.Time.UTC().Format(time.RFC3339), archived.ArchivedAt.Valid)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
