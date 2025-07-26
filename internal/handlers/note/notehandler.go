package notehandler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/sergeyreshetnyakov/notion/internal/domain/models"
	"github.com/sergeyreshetnyakov/notion/internal/lib/logger/sl"
)

type Handler struct {
	log   *slog.Logger
	notes Notes
}

type Notes interface {
	GetAll(ctx context.Context) (notes []models.Note, err error)
	Add(ctx context.Context, header string, content string) (err error)
	Edit(ctx context.Context, header string, content string, id int) (err error)
	Delete(ctx context.Context, id int) (err error)
}

func New(log *slog.Logger, notes Notes) Handler {
	return Handler{
		log:   log,
		notes: notes,
	}
}

func (h Handler) HandleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.GetAll)
	mux.HandleFunc("POST /", h.Add)
	mux.HandleFunc("PATCH /", h.Edit)
	mux.HandleFunc("DELETE /", h.Delete)
}

func (h Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "Note.GetAll"
	h.log.With(
		slog.String("op", op),
	)

	notes, err := h.notes.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to get notes: "+err.Error(), http.StatusInternalServerError)
		h.log.Info("Failed to get notes", sl.Err(err))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]models.Note{"notes": notes})
}

func (h Handler) Add(w http.ResponseWriter, r *http.Request) {
	const op = "Note.Add"
	h.log.With(
		slog.String("op", op),
	)

	var msg struct {
		Header  string `json:"header"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		h.log.Info("Failed to decode request body", sl.Err(err))
		return
	}

	if msg.Header == "" {
		err := errors.New("Header must contain any characters")
		http.Error(w, "Failed to add new note: "+err.Error(), http.StatusBadRequest)
		h.log.Info("Failed to add new note", sl.Err(err))
		return
	}

	if err := h.notes.Add(r.Context(), msg.Header, msg.Content); err != nil {
		http.Error(w, "Failed to add new note: "+err.Error(), http.StatusInternalServerError)
		h.log.Info("Failed to add new note", sl.Err(err))
		return
	}
}

func (h Handler) Edit(w http.ResponseWriter, r *http.Request) {
	const op = "Note.Edit"
	h.log.With(
		slog.String("op", op),
	)

	var msg struct {
		Header  string `json:"header"`
		Content string `json:"content"`
		Id      int    `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		h.log.Info("Failed to decode request body", sl.Err(err))
		return
	}

	if err := h.notes.Edit(r.Context(), msg.Header, msg.Content, msg.Id); err != nil {
		http.Error(w, "Failed to edit note: "+err.Error(), http.StatusInternalServerError)
		h.log.Info("Failed to edit note", sl.Err(err))
	}
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "Note.Delete"
	h.log.With(
		slog.String("op", op),
	)

	var msg struct {
		Id int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		h.log.Info("Failed to decode request body", sl.Err(err))
		return
	}

	err := h.notes.Delete(r.Context(), msg.Id)
	if err != nil {
		http.Error(w, "Failed to delete note: "+err.Error(), http.StatusInternalServerError)
		h.log.Info("Failed to delete note", sl.Err(err))
	}
}
