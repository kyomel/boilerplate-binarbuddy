package usecases

import (
	"boilerplate-sqlc/libs/helper"
	"boilerplate-sqlc/models"
	"boilerplate-sqlc/repositories"
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserUseCase interface {
	RegisterUser(ctx context.Context, request *models.RegisterRequest) (*models.User, error)
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

func (u *userUseCase) RegisterUser(ctx context.Context, request *models.RegisterRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	tx, err := u.db.Begin()
	if err != nil {
		log.Println("SQL error in UseCase RegisterUser => Open Transaction", err)
		return nil, err
	}

	userRegistered, err := u.userRepository.CheckRegistered(ctx, tx, request.Username)
	if err != nil {
		return nil, err
	}

	if userRegistered {
		return nil, errors.New("username already registered")
	}

	userHash, err := u.userRepository.GenerateUserHash(ctx, request.Password)
	if err != nil {
		return nil, err
	}

	userData, err := u.userRepository.RegisterUser(ctx, tx, &models.User{
		ID:       uuid.New().String(),
		Username: request.Username,
		Hash:     userHash,
	})

	if err != nil {
		return nil, err
	}

	helper.CommitOrRollback(tx, err)

	return userData, nil
}
