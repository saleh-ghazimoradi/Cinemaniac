package handlers

import (
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/dto"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/helper"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/repository"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/service"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/validator"
	"net/http"
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

	movie, err := m.movieService.CreateMovie(r.Context(), &payload)
	if err != nil {
		var valErr validator.ValidationError
		if errors.As(err, &valErr) {
			helper.FailedValidationResponse(w, r, valErr.Errors)
			return
		}
		helper.ServerErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	if err = helper.WriteJSON(w, http.StatusCreated, helper.Envelope{"movie": movie}, headers); err != nil {
		helper.ServerErrorResponse(w, r, err)
	}
}

func (m *MovieHandler) ShowMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ReadParams(r)
	if err != nil {
		helper.NotFoundResponse(w, r)
		return
	}

	movie, err := m.movieService.GetMovieById(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, r)
		default:
			helper.ServerErrorResponse(w, r, err)
		}
		return
	}

	if err = helper.WriteJSON(w, http.StatusOK, helper.Envelope{"movie": movie}, nil); err != nil {
		helper.ServerErrorResponse(w, r, err)
	}
}

func (m *MovieHandler) GetMoviesHandler(w http.ResponseWriter, r *http.Request) {
	movies, err := m.movieService.GetMovies(r.Context())
	if err != nil {
		helper.ServerErrorResponse(w, r, err)
		return
	}

	if err = helper.WriteJSON(w, http.StatusOK, helper.Envelope{"movie": movies}, nil); err != nil {
		helper.ServerErrorResponse(w, r, err)
	}
}

func (m *MovieHandler) UpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ReadParams(r)
	if err != nil {
		helper.NotFoundResponse(w, r)
		return
	}

	var input dto.UpdateMovie
	if err := helper.ReadJSON(w, r, &input); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	updatedMovie, err := m.movieService.UpdateMovie(r.Context(), id, &input)
	if err != nil {
		var valErr validator.ValidationError
		if errors.As(err, &valErr) {
			helper.FailedValidationResponse(w, r, valErr.Errors)
			return
		}
		
		switch {
		case errors.Is(err, repository.ErrEditConflict):
			helper.EditConflictResponse(w, r)
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, r)
		default:
			helper.ServerErrorResponse(w, r, err)
		}
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"movie": updatedMovie}, nil); err != nil {
		helper.ServerErrorResponse(w, r, err)
	}
}

func (m *MovieHandler) DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ReadParams(r)
	if err != nil {
		helper.NotFoundResponse(w, r)
		return
	}

	err = m.movieService.DeleteMovie(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, r)
		default:
			helper.ServerErrorResponse(w, r, err)
		}
		return
	}

	if err = helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "movie successfully deleted"}, nil); err != nil {
		helper.ServerErrorResponse(w, r, err)
	}
}

func NewMovieHandler(movieService service.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}
