package usecases

import (
	"boilerplate-sqlc/libs/helper"
	"boilerplate-sqlc/models"
	"boilerplate-sqlc/repositories"
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type AuthorUseCase interface {
	CreateAuthor(ctx context.Context, req *models.AuthorReq) (*models.AuthorResp, error)
}

type authorUseCase struct {
	contextTimeout   time.Duration
	authorRepository repositories.AuthorRepository
	db               *sqlx.DB
}

func NewAuthorUseCase(contextTimeout time.Duration, authorRepository repositories.AuthorRepository, db *sqlx.DB) AuthorUseCase {
	return &authorUseCase{
		contextTimeout,
		authorRepository,
		db,
	}
}

func (u *authorUseCase) CreateAuthor(ctx context.Context, req *models.AuthorReq) (*models.AuthorResp, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	tx, err := u.db.Begin()
	if err != nil {
		log.Println("SQL error in UseCase CreateAuthor => Open Transaction", err)
		return nil, err
	}

	author, err := u.authorRepository.CreateAuthor(ctx, tx, req)
	if err != nil {
		return nil, err
	}

	helper.CommitOrRollback(tx, err)

	return author, nil
}
