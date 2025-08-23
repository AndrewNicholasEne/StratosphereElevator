package stacks

import (
	"context"
	"errors"
	"log/slog"

	"github.com/AndrewNicholasEne/StratosphereElevator/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Service struct {
	s Store
	l *slog.Logger
}

func New(s Store, l *slog.Logger) *Service {
	if l == nil {
		l = slog.Default()
	}
	return &Service{s: s, l: l}
}

type CreateInput struct {
	Name string
	Slug *string
}

func (sv *Service) Create(ctx context.Context, in CreateInput) (db.Stack, error) {
	name, slug, err := normalizeCreateInput(in)
	if err != nil {
		return db.Stack{}, ErrInvalidInput
	}

	stack, err := sv.s.CreateStack(ctx, db.CreateStackParams{
		ID:   uuid.New(),
		Name: name,
		Slug: slug,
	})

	if err != nil {
		var pgerr *pgconn.PgError
		if isUniqueViolation(pgerr, "") {
			return db.Stack{}, ErrConflict
		}
		return db.Stack{}, err
	}

	return stack, nil
}

func (sv *Service) GetBySlug(ctx context.Context, slug string) (db.Stack, error) {
	stack, err := sv.s.GetStackBySlug(ctx, slugify(slug))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Stack{}, ErrNotFound
		}
		return db.Stack{}, err
	}
	return stack, nil
}

func (sv *Service) List(ctx context.Context, params db.ListStacksParams) ([]db.Stack, error) {
	stacks, err := sv.s.ListStacks(ctx, params)
	if err != nil {
		return []db.Stack{}, err
	}
	if stacks == nil {
		stacks = make([]db.Stack, 0)
	}

	return stacks, nil
}

func (sv *Service) Archive(ctx context.Context, id uuid.UUID) error {
	rows, err := sv.s.ArchiveStack(ctx, id)
	if err != nil {
		return err
	}
	if rows == 1 {
		return nil
	}

	already, err := sv.s.StackArchivedStatus(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if already {
		return ErrAlreadyArchived
	}
	return ErrNotFound
}
