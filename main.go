package main

import (
	"integration-test-demo/book"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	bookRepo := book.NewRepo()
	bookService := book.NewService(bookRepo)
	bookHandler := book.NewHandler(bookService)

	router.Route("/book", bookHandler.API())

	log.Println("Server starting on port 8848...")
	err := http.ListenAndServe(":8848", router)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
