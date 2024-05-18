package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Hash     string `json:"-"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSession struct {
	JWTToken string `json:"jwt_token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
