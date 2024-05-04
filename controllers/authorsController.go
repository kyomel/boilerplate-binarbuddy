package controllers

import (
	"boilerplate-sqlc/libs/routers"
	"boilerplate-sqlc/libs/utils"
	"boilerplate-sqlc/models"
	"boilerplate-sqlc/usecases"
	"encoding/json"
	"net/http"
)

type AuthorsController interface {
	CreateAuthor(w http.ResponseWriter, r *http.Request)
}

type authorController struct {
	authorUC usecases.AuthorUseCase
	r        routers.Resultset
}

func NewAuthorsController(authorUC usecases.AuthorUseCase, r routers.Resultset) AuthorsController {
	return &authorController{
		authorUC,
		r,
	}
}

func (c *authorController) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var (
		request     models.AuthorReq
		successResp models.SuccessResponse
		errResp     models.ErrorResponse
	)

	binding := json.NewDecoder(r.Body).Decode(&request)
	if binding != nil {
		errResp.Status = http.StatusBadRequest
		errResp.Error = utils.BAD_REQUEST
		errResp.Message = http.StatusText(http.StatusBadRequest)

		c.r.ResponsWithError(w, http.StatusBadRequest, utils.BAD_REQUEST)
		return
	}

	res, err := c.authorUC.CreateAuthor(r.Context(), &request)
	if err != nil {
		errResp.Status = http.StatusInternalServerError
		errResp.Error = utils.INTERNAL_SERVER_ERROR
		errResp.Message = http.StatusText(http.StatusInternalServerError)

		c.r.ResponsWithError(w, http.StatusInternalServerError, utils.INTERNAL_SERVER_ERROR)
		return
	}

	successResp.Data = res
	successResp.Status = http.StatusCreated
	successResp.Message = "User created successfully"
	c.r.ResponsWithJSON(w, http.StatusCreated, successResp)
}
