package usecases

import (
	"boilerplate-sqlc/models"
	"boilerplate-sqlc/repositories"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserUseCase interface {
	RegisterUser(ctx context.Context, userData *models.User) (*models.User, error)
}

type userUseCase struct {
	contextTimeout time.Duration
	userRepository repositories.UserRepository
	db             *sqlx.DB
}

func NewUserUseCase(contextTimeout time.Duration, userRepository repositories.UserRepository, db *sqlx.DB) UserUseCase {
	return &userUseCase{
		contextTimeout,
		userRepository,
		db,
	}
}

func (u *userUseCase) RegisterUser(ctx context.Context, userData *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return nil, nil
}
