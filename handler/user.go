package handler

import (
	"fmt"
	"net/http"
	"yurafund/auth"
	"yurafund/helper"
	"yurafund/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authservice auth.Service
}

func NewUserHandler(userService user.Service, authservice auth.Service) *userHandler {
	return &userHandler{userService, authservice}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		fmt.Println(err)
		errors := helper.FormatValidationError(err)
		responseError := helper.APIResponse("Failed to process request", http.StatusUnprocessableEntity, "error", errors)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}
	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		fmt.Println(err)
		responseError := helper.APIResponse("Failed to process request", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}
	// set token jwt
	jwtToken, err := h.authservice.GenerateToken(newUser.ID)
	if err != nil {
		fmt.Println(err)
		responseError := helper.APIResponse("Failed to process request", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}

	formatter := user.FormatUser(newUser, jwtToken)
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
		errorsMessage := gin.H{"errors": err.Error()}
		responseError := helper.APIResponse("Error Happend during login", http.StatusBadRequest, "error", errorsMessage)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}

	// set token jwt
	jwtToken, err := h.authservice.GenerateToken(loginUser.ID)
	if err != nil {
		fmt.Println(err)
		responseError := helper.APIResponse("Failed to process request", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}

	formatter := user.FormatUser(loginUser, jwtToken)
	response := helper.APIResponse("Successfuly loggedin", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		responseError := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errors)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errors := gin.H{"errors": "Server error"}
		responseError := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errors)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}

	metaMessage := "Email has been registered"
	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	data := gin.H{"is_available": isEmailAvailable}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfuly uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
