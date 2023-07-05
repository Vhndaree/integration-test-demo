package book

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type IService interface {
	GetAll() []*Book
	Post(book *Book) error
	Get(id string) (*Book, error)
	Update(id string, book *Book) error
	Delete(id string)
}

type Handler struct {
	service IService
}

func NewHandler(service IService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) API() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", h.getAll())
		r.Post("/", h.post())
		r.With(h.validate()).Get("/{id}", h.get())
		r.With(h.validate()).Put("/{id}", h.update())
		r.With(h.validate()).Delete("/{id}", h.delete())
	}
}

func (h *Handler) getAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books := h.service.GetAll()
		if books == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if err := json.NewEncoder(w).Encode(&books); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload Book
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] invalid body (%s)", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = h.service.Post(&payload)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] (%s)", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		book, err := h.service.Get(id)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] (%s)", err))

			statusCode := http.StatusInternalServerError
			if errors.Is(err, ErrNotFound) {
				statusCode = http.StatusNotFound
			}

			w.WriteHeader(statusCode)
			return
		}

		if err := json.NewEncoder(w).Encode(&book); err != nil {
			fmt.Println(fmt.Sprintf("[ERROR] (%s)", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var payload Book
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] invalid body (%s)", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = h.service.Update(id, &payload)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR] (%s)", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		h.service.Delete(id)
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) validate() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			id := chi.URLParam(req, "id")
			_, err := uuid.Parse(id)
			if err != nil {
				log.Println(fmt.Sprintf("[ERROR] invalid id (%s)", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}
