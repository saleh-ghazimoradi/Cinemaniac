package handlers

import (
	"fmt"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/domain"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/helper"
	"net/http"
	"time"
)

type MovieHandler struct{}

func (m *MovieHandler) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create movie")
}

func (m *MovieHandler) ShowMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ReadParams(r)
	if err != nil {
		helper.NotFoundResponse(w, r)
		return
	}

	movie := domain.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "The Black List",
		Runtime:   102,
		Genres:    []string{"drama", "criminal", "FBI"},
		Version:   1,
	}

	if err = helper.WriteJSON(w, http.StatusOK, helper.Envelope{"movie": movie}, nil); err != nil {
		helper.ServerErrorResponse(w, r, err)
	}
}

func NewMovieHandler() *MovieHandler {
	return &MovieHandler{}
}
