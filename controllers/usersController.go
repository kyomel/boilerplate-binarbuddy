package controllers

import (
	"boilerplate-sqlc/libs/routers"
	"boilerplate-sqlc/libs/utils"
	"boilerplate-sqlc/models"
	"boilerplate-sqlc/usecases"
	"encoding/json"
	"net/http"
)

type UsersController interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	userUC usecases.UserUseCase
	r      routers.Resultset
}

func NewUsersController(userUC usecases.UserUseCase, r routers.Resultset) UsersController {
	return &userController{
		userUC,
		r,
	}
}

func (c *userController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var (
		request     models.RegisterRequest
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

	res, err := c.userUC.RegisterUser(r.Context(), &request)
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
