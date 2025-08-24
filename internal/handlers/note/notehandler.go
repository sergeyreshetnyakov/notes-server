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
	Add(ctx context.Context, header string, content string) (id int64, err error)
	Edit(ctx context.Context, header string, content string, id int64) (err error)
	Delete(ctx context.Context, id int64) (err error)
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

// GetAll godoc
//
//	@Summary		Get all notes
//	@Description	Returns a list of notes
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number"
//	@Param			results	query		int	false	"Results per page"
//	@Success		200		{object}	[]models.Note
//	@Failure		404		{string}	string	"page not found"
//	@Failure		500		{string}	string	"internal server error"
//	@Router			/ [get]
func (h Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "Note.GetAll"
	h.log.With(
		slog.String("op", op),
	)

	notes, err := h.notes.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to get notes: "+err.Error(), http.StatusInternalServerError)
		h.log.Error("Failed to get notes", sl.Err(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// AddNote godoc
//
//	@Summary		Add note
//	@Description	Adds a new note
//	@Accept			json
//	@Produce		json
//	@Param			header	body	string	true	"Notes header"
//	@Param			content	body	string	true	"Notes content"
//	@Success		200
//	@Failure		400	{string}	string	"bad request body"
//	@Failure		500	{string}	string	"internal server error"
//	@Router			/ [post]
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
		h.log.Debug("Failed to decode request body", sl.Err(err))
		return
	}

	if msg.Header == "" {
		err := errors.New("Header must contain any characters")
		http.Error(w, "Failed to add new note: "+err.Error(), http.StatusBadRequest)
		h.log.Debug("Failed to add new note", sl.Err(err))
		return
	}

	id, err := h.notes.Add(r.Context(), msg.Header, msg.Content)
	if err != nil {
		http.Error(w, "Failed to add new note: "+err.Error(), http.StatusInternalServerError)
		h.log.Error("Failed to add new note", sl.Err(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// EditNote godoc
//
//	@Summary		Edit note
//	@Description	Edits a note
//	@Accept			json
//	@Produce		json
//	@Param			note	body	models.Note	true	"Notes body"
//	@Success		200
//	@Failure		400	{string}	string	"bad request body"
//	@Failure		404	{string}	string	"note not found"
//	@Failure		500	{string}	string	"internal server error"
//	@Router			/ [patch]
func (h Handler) Edit(w http.ResponseWriter, r *http.Request) {
	const op = "Note.Edit"
	h.log.With(
		slog.String("op", op),
	)

	var msg struct {
		Header  string `json:"header"`
		Content string `json:"content"`
		Id      int64  `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		h.log.Debug("Failed to decode request body", sl.Err(err))
		return
	}

	if err := h.notes.Edit(r.Context(), msg.Header, msg.Content, msg.Id); err != nil {
		if errors.Is(err, errors.New("note not found")) {
			http.Error(w, "Failed to edit note:"+err.Error(), http.StatusNotFound)
			h.log.Debug("Failed to edit note", sl.Err(err))
		} else {
			http.Error(w, "Failed to edit note: "+err.Error(), http.StatusInternalServerError)
			h.log.Error("Failed to edit note", sl.Err(err))
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteNote godoc
//
//	@Summary		Delete note
//	@Description	Deletes a note
//	@Accept			json
//	@Produce		json
//	@Param			id	body	int	true	"Notes id"
//	@Success		200
//	@Failure		400	{string}	string	"bad request body"
//	@Failure		404	{string}	string	"note not found"
//	@Failure		500	{string}	string	"internal server error"
//	@Router			/ [delete]
func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "Note.Delete"
	h.log.With(
		slog.String("op", op),
	)

	var msg struct {
		Id int64 `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		h.log.Debug("Failed to decode request body", sl.Err(err))
		return
	}

	err := h.notes.Delete(r.Context(), msg.Id)
	if err != nil {
		if errors.Is(err, errors.New("note not found")) {
			http.Error(w, "Failed to delete note:"+err.Error(), http.StatusNotFound)
			h.log.Debug("Failed to delete note", sl.Err(err))
		} else {
			http.Error(w, "Failed to delete note: "+err.Error(), http.StatusInternalServerError)
			h.log.Error("Failed to delete note", sl.Err(err))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
