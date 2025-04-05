package handlers

import (
	"errors"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/dto"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/helper"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/repository"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/service"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/validator"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func (u *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload *dto.User

	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	user, err := u.userService.CreateUser(r.Context(), payload)
	if err != nil {
		var valErr validator.ValidationError
		if errors.As(err, &valErr) {
			helper.FailedValidationResponse(w, r, valErr.Errors)
			return
		}
		switch {
		case errors.Is(err, repository.ErrDuplicateEmail):
			helper.ErrorResponse(w, r, http.StatusConflict, "a user with this email address already exists")
		default:
			helper.ServerErrorResponse(w, r, err)
		}
		return
	}
	if err = helper.WriteJSON(w, http.StatusCreated, helper.Envelope{"user": user}, nil); err != nil {
		helper.ServerErrorResponse(w, r, err)
	}
}

func (u *UserHandler) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload *dto.ActivateUserRequest

	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	user, err := u.userService.ActivateUser(r.Context(), payload)
	if err != nil {
		var valErr validator.ValidationError
		switch {
		case errors.As(err, &valErr):
			helper.FailedValidationResponse(w, r, valErr.Errors)
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.ErrorResponse(w, r, http.StatusUnprocessableEntity, "invalid or expired activation token")
		case errors.Is(err, repository.ErrEditConflict):
			helper.EditConflictResponse(w, r)
		default:
			helper.ServerErrorResponse(w, r, err)
		}
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"user": user}, nil); err != nil {
		helper.ServerErrorResponse(w, r, err)
	}
}

func (u *UserHandler) CreateAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload *dto.Token
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	token, err := u.userService.CreateAuthenticationToken(r.Context(), payload)
	if err != nil {
		var valErr validator.ValidationError
		switch {
		case errors.As(err, &valErr):
			helper.FailedValidationResponse(w, r, valErr.Errors)
		case err.Error() == "invalid credentials":
			helper.InvalidCredentialsResponse(w, r)
		default:
			helper.ServerErrorResponse(w, r, err)
		}
		return
	}

	if err := helper.WriteJSON(w, http.StatusCreated, helper.Envelope{"authentication_token": token}, nil); err != nil {
		helper.ServerErrorResponse(w, r, err)
	}
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
