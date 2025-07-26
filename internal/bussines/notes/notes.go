package notes

import (
	"context"
	"errors"

	"github.com/sergeyreshetnyakov/notion/internal/domain/models"
)

type Storage interface {
	GetAll(ctx context.Context) (notes []models.Note, err error)
	GetById(ctx context.Context, id int) (note models.Note, err error)
	Add(ctx context.Context, header string, content string) (err error)
	Edit(ctx context.Context, header string, content string, id int) (err error)
	Delete(ctx context.Context, id int) (err error)
}

type Notes struct {
	storage Storage
}

func New(storage Storage) Notes {
	return Notes{storage}
}

func (n Notes) GetAll(ctx context.Context) (notes []models.Note, err error) {
	notes, err = n.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (n Notes) Add(ctx context.Context, header string, content string) (err error) {
	err = n.storage.Add(ctx, header, content)
	return err
}

func (n Notes) Edit(ctx context.Context, header string, content string, id int) (err error) {
	note, err := n.storage.GetById(ctx, id)
	if err != nil {
		return err
	}

	if header == "" {
		header = note.Header
	}

	if content == "" {
		content = note.Content
	}

	if header == note.Header && content == note.Content {
		return errors.New("nothing to change")
	}

	err = n.storage.Edit(ctx, header, content, id)
	return err
}

func (n Notes) Delete(ctx context.Context, id int) (err error) {
	err = n.storage.Delete(ctx, id)
	return err
}
