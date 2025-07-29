package notestorage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sergeyreshetnyakov/notion/internal/domain/models"
)

type Storage struct {
	db *sql.DB
}

var ErrNoteNotFound = errors.New("note not found")

type shutdownFunc func() error

func New(storagePath string, log *slog.Logger) (*Storage, shutdownFunc) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		panic(err)
	}
	log.Info("DB is connected")

	return &Storage{db: db}, func() error {
		err := db.Close()
		return err
	}
}

func (s *Storage) GetAll(ctx context.Context) (notes []models.Note, err error) {
	stmt, err := s.db.Prepare("SELECT header, content, id FROM notes")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	for rows.Next() {
		var header, content string
		var id int64
		if err := rows.Scan(&header, &content, &id); err != nil {
			return nil, err
		}
		notes = append(notes, models.Note{
			Header:  header,
			Content: content,
			Id:      id,
		})
	}

	return notes, nil
}

func (s *Storage) GetById(ctx context.Context, id int64) (note models.Note, err error) {
	stmt, err := s.db.Prepare("SELECT header, content, id FROM notes WHERE id = ?")
	if err != nil {
		return models.Note{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	if err := row.Scan(&note.Header, &note.Content, &note.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Note{}, ErrNoteNotFound
		}
		return models.Note{}, err
	}

	return note, nil
}

func (s *Storage) Add(ctx context.Context, header string, content string) (id int64, err error) {
	stmt, err := s.db.Prepare("INSERT INTO notes(header, content) VALUES(?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, header, content)
	if err != nil {
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func (s *Storage) Edit(ctx context.Context, header string, content string, id int64) (err error) {
	stmt, err := s.db.Prepare("UPDATE notes SET header = ?, content = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, header, content, id)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); rows == 0 {
		if err != nil {
			return err
		}

		return ErrNoteNotFound
	}

	return err
}

func (s *Storage) Delete(ctx context.Context, id int64) (err error) {
	stmt, err := s.db.Prepare("DELETE FROM notes WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if rows, err := res.RowsAffected(); rows == 0 {
		if err != nil {
			return err
		}
		return ErrNoteNotFound
	}
	return err
}
