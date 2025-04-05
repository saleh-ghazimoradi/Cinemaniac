package service

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/domain"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/dto"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/repository"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/transaction"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/validator"
	"github.com/saleh-ghazimoradi/Cinemaniac/pkg/notification"
	"github.com/saleh-ghazimoradi/Cinemaniac/slg"
	"github.com/saleh-ghazimoradi/Cinemaniac/utils"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, input *dto.User) (*domain.User, error)
	ActivateUser(ctx context.Context, input *dto.ActivateUserRequest) (*domain.User, error)
}

type userService struct {
	userRepository  repository.UserRepository
	txService       transaction.TxService
	notification    notification.Mailer
	tokenRepository repository.TokenRepository
}

func (u *userService) CreateUser(ctx context.Context, input *dto.User) (*domain.User, error) {
	v := validator.New()

	domain.ValidateEmail(v, input.Email)
	domain.ValidatePasswordPlaintext(v, input.Password)

	if err := v.GetValidationError(); err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  domain.Password{},
		Activated: false,
	}

	if err := user.Password.Set(input.Password); err != nil {
		return nil, err
	}

	var token *domain.Token

	if err := u.txService.WithTx(ctx, func(tx *sql.Tx) error {
		txUserRepo := u.userRepository.WithTx(ctx, tx)
		if err := txUserRepo.CreateUser(ctx, user); err != nil {
			return err
		}

		token = utils.GenerateToken(user.ID, 72*time.Hour, domain.ScopeActivation)

		txTokenRepo := u.tokenRepository.WithTx(ctx, tx)
		return txTokenRepo.Insert(ctx, token)
	}); err != nil {
		return nil, err
	}

	background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userId":          user.ID,
		}

		err := u.notification.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			slg.Logger.Error(err.Error())
		}
	})

	return user, nil
}

func (u *userService) ActivateUser(ctx context.Context, input *dto.ActivateUserRequest) (*domain.User, error) {
	v := validator.New()

	if domain.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		return nil, v.GetValidationError()
	}

	user, err := u.userRepository.GetForToken(ctx, domain.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		return nil, err
	}

	user.Activated = true

	err = u.txService.WithTx(ctx, func(tx *sql.Tx) error {
		txUserRepo := u.userRepository.WithTx(ctx, tx)
		txTokenRepo := u.tokenRepository.WithTx(ctx, tx)

		if err := txUserRepo.UpdateUser(ctx, user); err != nil {
			return err
		}

		if err := txTokenRepo.DeleteAllForUser(ctx, domain.ScopeActivation, user.ID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserService(userRepository repository.UserRepository, txService transaction.TxService, notification notification.Mailer, tokenRepository repository.TokenRepository) UserService {
	return &userService{
		userRepository:  userRepository,
		txService:       txService,
		notification:    notification,
		tokenRepository: tokenRepository,
	}
}
