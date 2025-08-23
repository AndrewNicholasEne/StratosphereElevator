package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/AndrewNicholasEne/StratosphereElevator/internal/db"
	"github.com/AndrewNicholasEne/StratosphereElevator/internal/stacks"
	"github.com/google/uuid"
)

type StacksHTTP struct{ svc *stacks.Service }

func NewStacksHTTP(svc *stacks.Service) *StacksHTTP { return &StacksHTTP{svc: svc} }

func (h *StacksHTTP) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /stacks", h.List)
	mux.HandleFunc("GET /stacks/{slug}", h.GetBySlug)
	mux.HandleFunc("POST /stacks", h.Create)
	mux.HandleFunc("POST /stacks/{id}/archive", h.Archive)
}

type stackResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Slug       string     `json:"slug"`
	CreatedAt  time.Time  `json:"created_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty"`
}

func toStackResponse(s db.Stack) stackResponse {
	var at *time.Time
	if s.ArchivedAt.Valid {
		t := s.ArchivedAt.Time
		at = &t
	}
	return stackResponse{ID: s.ID, Name: s.Name, Slug: s.Slug, CreatedAt: s.CreatedAt, ArchivedAt: at}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func (h *StacksHTTP) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context(), db.ListStacksParams{Column1: false})
	if err != nil {
		writeErr(w, 500, "internal")
		return
	}
	out := make([]stackResponse, 0, len(items))
	for _, s := range items {
		out = append(out, toStackResponse(s))
	}
	writeJSON(w, 200, out)
}
func (h *StacksHTTP) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	s, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, stacks.ErrNotFound) {
			writeErr(w, 404, "not found")
			return
		}
		writeErr(w, 500, "internal")
		return
	}
	writeJSON(w, 200, toStackResponse(s))
}

func (h *StacksHTTP) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string
		Slug *string
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "bad json")
		return
	}
	st, err := h.svc.Create(r.Context(), stacks.CreateInput{Name: req.Name, Slug: req.Slug})
	if err != nil {
		switch {
		case errors.Is(err, stacks.ErrInvalidInput):
			writeErr(w, 400, "invalid input")
		case errors.Is(err, stacks.ErrConflict):
			writeErr(w, 409, "slug conflict")
		default:
			writeErr(w, 500, "internal")
		}
		return
	}
	w.Header().Set("Location", "/stacks/"+st.Slug)
	writeJSON(w, 201, toStackResponse(st))
}

func (h *StacksHTTP) Archive(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeErr(w, 400, "invalid id")
		return
	}
	if err := h.svc.Archive(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, stacks.ErrAlreadyArchived):
			writeErr(w, 409, "already archived")
		case errors.Is(err, stacks.ErrNotFound):
			writeErr(w, 404, "not found")
		default:
			writeErr(w, 500, "internal")
		}
		return
	}
	w.WriteHeader(204)
}
