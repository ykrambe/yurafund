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
	// set token jwt
	formatter := user.FormatUser(newUser, "tokentok")
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		responseError := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errors) // 422
		c.JSON(http.StatusBadRequest, responseError)
		return
	}

	loginUser, err := h.userService.Login(input)
	if err != nil {
		errosMessage := gin.H{"errors": err.Error()}
		responseError := helper.APIResponse("Error Happend during login", http.StatusBadRequest, "error", errosMessage)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}

	formatter := user.FormatUser(loginUser, "tokentok")
	response := helper.APIResponse("Successfuly loggedin", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}
