package repositories

import (
	"boilerplate-sqlc/models"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, userData *models.User) (*models.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) RegisterUser(ctx context.Context, userData *models.User) (*models.User, error) {
	var user models.User
	query := `
		INSERT INTO users(id, username) 
		VALUES ($1, $2) 
		RETURNING id, username
	`
	row := r.db.QueryRowContext(ctx, query, userData.Username, userData.Hash)
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		log.Println("SQL error on RegisterUser => Execute Query and Scan", err)
		return &user, err
	}

	return &user, nil
}
