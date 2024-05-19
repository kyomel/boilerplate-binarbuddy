package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	c "boilerplate-sqlc/controllers"
	databases "boilerplate-sqlc/db"
	router "boilerplate-sqlc/libs/routers"
	"boilerplate-sqlc/repositories"
	"boilerplate-sqlc/usecases"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

var (
	httpRouter router.Router          = router.NewChiRouter()
	dbRepoConn databases.DatabaseRepo = databases.NewPostgresRepo()
)

func main() {
	// Load envornment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Read environment variables
	port := os.Getenv("PORT")
	appName := os.Getenv("APP_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("SSL_MODE")

	timeoutContext := time.Duration(5 * time.Second)
	httpResult := router.NewResultset()

	portDB, _ := strconv.Atoi(dbPort)
	portConnect, _ := strconv.Atoi(port)

	// Auth Settings
	secret := "AES256Key-32Characters1234567890"
	signKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	db, err := dbRepoConn.Connect(dbHost, portDB, dbUser, dbPassword, dbName, sslMode)
	if err != nil {
		log.Fatal(err)
	}

	// Call Repositories
	authorRepo := repositories.NewAuthorRepository(db)
	userRepo := repositories.NewUserRepository(db, secret, 1, 64*1024, 4, 32, signKey, 60*time.Second)

	// Call UseCases
	authorUC := usecases.NewAuthorUseCase(timeoutContext, authorRepo, db)
	userUC := usecases.NewUserUseCase(timeoutContext, userRepo, db)

	// Call Controllers
	author := c.NewAuthorsController(authorUC, httpResult)
	user := c.NewUsersController(userUC, httpResult)

	// Define a simple GET route
	httpRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Welcome to %s!", appName)
	})

	httpRouter.Post("/authors", author.CreateAuthor)
	httpRouter.Post("/users/register", user.RegisterUser)

	// Start the server
	fmt.Printf("Starting %s on port %d\n", appName, portConnect)

	httpRouter.Run(portConnect, appName)
}
