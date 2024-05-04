package models

type AuthorReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AuthorResp struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
