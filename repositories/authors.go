package repositories

import (
	"boilerplate-sqlc/models"
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type AuthorRepository interface {
	CreateAuthor(ctx context.Context, tx *sql.Tx, req *models.AuthorReq) (*models.AuthorResp, error)
}

type authorRepository struct {
	db *sqlx.DB
}

func NewAuthorRepository(db *sqlx.DB) AuthorRepository {
	return &authorRepository{
		db,
	}
}

func (r *authorRepository) CreateAuthor(ctx context.Context, tx *sql.Tx, req *models.AuthorReq) (*models.AuthorResp, error) {
	var author models.AuthorResp
	query := `
		INSERT INTO author(name, email) 
		VALUES ($1, $2) 
		RETURNING id, name, email
	`
	row := tx.QueryRowContext(ctx, query, req.Name, req.Email)
	err := row.Scan(&author.ID, &author.Name, &author.Email)
	if err != nil {
		log.Println("SQL error on CreateAuthor => Execute Query and Scan", err)
		return nil, err
	}

	return &author, nil
}
