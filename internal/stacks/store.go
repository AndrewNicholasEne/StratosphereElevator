package stacks

import (
	"context"

	"github.com/AndrewNicholasEne/StratosphereElevator/internal/db"
	"github.com/google/uuid"
)

type Store interface {
	CreateStack(context.Context, db.CreateStackParams) (db.Stack, error)
	GetStackBySlug(context.Context, string) (db.Stack, error)
	ListStacks(context.Context, db.ListStacksParams) ([]db.Stack, error)
	ArchiveStack(ctx context.Context, id uuid.UUID) (int64, error)
	StackArchivedStatus(ctx context.Context, id uuid.UUID) (bool, error)
}
