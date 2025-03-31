package handlers

import (
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/domain"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/dto"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/helper"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/service"
	"net/http"
	"time"
)

type MovieHandler struct {
	movieService service.MovieService
}

func (m *MovieHandler) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var payload dto.Movie

	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	movie, validationErrors, err := m.movieService.CreateMovie(r.Context(), &payload)
	if err != nil {
		if validationErrors != nil {
			helper.FailedValidationResponse(w, r, validationErrors)
			return
		}

		helper.ServerErrorResponse(w, r, err)
		return
	}

	if err = helper.WriteJSON(w, http.StatusOK, helper.Envelope{"movie": movie}, nil); err != nil {
		helper.ServerErrorResponse(w, r, err)
		return
	}
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

func NewMovieHandler(movieService service.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}
