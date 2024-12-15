package handler

import (
	"net/http"
	"yurafund/helper"
	"yurafund/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {

		errors := helper.FormatValidationError(err)
		responseError := helper.APIResponse("Failed to process request", http.StatusUnprocessableEntity, "error", errors)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}
	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		responseError := helper.APIResponse("Failed to process request", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}

	formatter := user.FormatUser(newUser, "tokentok")
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}
