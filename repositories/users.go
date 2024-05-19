package repositories

import (
	"boilerplate-sqlc/models"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/argon2"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, tx *sql.Tx, userData *models.User) (*models.User, error)
	CheckRegistered(ctx context.Context, tx *sql.Tx, username string) (bool, error)
	GenerateUserHash(ctx context.Context, password string) (hash string, err error)
}

type userRepository struct {
	db        *sqlx.DB
	gcm       cipher.AEAD
	time      uint32
	memory    uint32
	threads   uint8
	keylen    uint32
	signKey   *rsa.PrivateKey
	accessExp time.Duration
}

func NewUserRepository(db *sqlx.DB, secret string, time uint32, memory uint32, threads uint8, keylen uint32, signKey *rsa.PrivateKey, accessExp time.Duration) UserRepository {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}

	return &userRepository{
		db,
		gcm,
		time,
		memory,
		threads,
		keylen,
		signKey,
		accessExp,
	}
}

const (
	cryptForm = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
)

func (r *userRepository) RegisterUser(ctx context.Context, tx *sql.Tx, userData *models.User) (*models.User, error) {
	var user models.User
	query := `
		INSERT INTO users(id, username, hash) 
		VALUES ($1, $2, $3) 
		RETURNING id, username, hash
	`
	row := tx.QueryRowContext(ctx, query, userData.ID, userData.Username, userData.Hash)
	err := row.Scan(&user.ID, &user.Username, &user.Hash)
	if err != nil {
		log.Println("SQL error on RegisterUser => Execute Query and Scan", err)
		return &user, err
	}

	return &user, nil
}

func (r *userRepository) CheckRegistered(ctx context.Context, tx *sql.Tx, username string) (bool, error) {
	var exists bool
	query := `
	SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE username = $1
	)`

	err := tx.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		log.Println("SQL error on CheckRegistered => Execute Query and Scan", err)
		return false, err
	}

	return exists, nil
}

func (r *userRepository) GenerateUserHash(ctx context.Context, password string) (hash string, err error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	argonHash := argon2.IDKey([]byte(password), salt, r.time, r.memory, r.threads, r.keylen)

	b64hash := r.encrypt(ctx, argonHash)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)

	encodedHash := fmt.Sprintf(cryptForm, argon2.Version, r.memory, r.time, r.threads, b64Salt, b64hash)

	return encodedHash, nil
}

func (r *userRepository) encrypt(ctx context.Context, text []byte) string {
	nonce := make([]byte, r.gcm.NonceSize())

	ciphertext := r.gcm.Seal(nonce, nonce, text, nil)

	return base64.StdEncoding.EncodeToString(ciphertext)
}
