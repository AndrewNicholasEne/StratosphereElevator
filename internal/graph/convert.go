package graph

import (
	"time"

	"github.com/AndrewNicholasEne/StratosphereElevator/internal/db"
	"github.com/AndrewNicholasEne/StratosphereElevator/internal/graph/model"
)

func toGQL(s db.Stack) model.Stack {
	var at *time.Time
	if s.ArchivedAt.Valid {
		t := s.ArchivedAt.Time
		at = &t
	}
	return model.Stack{
		ID:         s.ID.String(),
		Name:       s.Name,
		Slug:       s.Slug,
		CreatedAt:  s.CreatedAt,
		ArchivedAt: at,
	}
}
