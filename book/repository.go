package book

import (
	"fmt"
)

type Repository struct {
	books map[string]*Book
}

func NewRepo() *Repository {
	return &Repository{books: make(map[string]*Book)}
}

func (r *Repository) GetAll() (books []*Book) {
	for _, book := range r.books {
		books = append(books, book)
	}

	return books
}

func (r *Repository) Get(id string) *Book {
	if book, ok := r.books[id]; ok {
		return book
	}

	return nil
}

func (r *Repository) Post(book *Book) {
	r.books[book.ID.String()] = book
}

func (r *Repository) Put(book *Book) error {
	book, ok := r.books[book.ID.String()]
	if !ok {
		return fmt.Errorf("[ERROR] Resource not exist in the store")
	}

	r.books[book.ID.String()] = book
	return nil
}

func (r *Repository) Delete(id string) {
	delete(r.books, id)
}
