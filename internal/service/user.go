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
)

type UserService interface {
	CreateUser(ctx context.Context, input *dto.User) (*domain.User, error)
}

type userService struct {
	userRepository repository.UserRepository
	txService      transaction.TxService
	notification   notification.Mailer
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

	if err := u.txService.WithTx(ctx, func(tx *sql.Tx) error {
		txUserRepo := u.userRepository.WithTx(ctx, tx)
		return txUserRepo.CreateUser(ctx, user)
	}); err != nil {
		return nil, err
	}

	background(func() {
		err := u.notification.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			slg.Logger.Error(err.Error())
		}
	})

	return user, nil
}

func NewUserService(userRepository repository.UserRepository, txService transaction.TxService, notification notification.Mailer) UserService {
	return &userService{
		userRepository: userRepository,
		txService:      txService,
		notification:   notification,
	}
}
