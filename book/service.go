package book

import (
	"fmt"

	"github.com/google/uuid"
)

type IRepository interface {
	GetAll() (books []*Book)
	Post(book *Book)
	Get(id string) *Book
	Put(id string, book *Book) error
	Delete(id string)
}

type Service struct {
	repo IRepository
}

func NewService(repo IRepository) *Service {
	return &Service{repo: repo}
}

// GetAll() ([]*Book, error)
// 	Post(book *Book) error
// 	Get(id string) (*Book, error)
// 	Update(id string, book *Book) error
// 	Delete(id string) error

func (s *Service) GetAll() []*Book {
	return s.repo.GetAll()
}

func (s *Service) Post(book *Book) (*Book, error) {
	if book.Title == "" {
		return nil, fmt.Errorf("invalid data")
	}

	book.ID = uuid.New()
	s.repo.Post(book)
	return book, nil
}
func (s *Service) Get(id string) (*Book, error) {
	book := s.repo.Get(id)
	if book == nil {
		return nil, ErrNotFound
	}

	return book, nil
}

func (s *Service) Update(id string, book *Book) error {
	oldBook := s.repo.Get(id)
	if oldBook == nil {
		return ErrNotFound
	}

	return s.repo.Put(id, book)
}

func (s *Service) Delete(id string) {
	s.repo.Delete(id)
}
