package handler

import (
	"net/http"
	"yurafund/helper"
	"yurafund/transaction"
	"yurafund/user"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed get transaction of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	listTransaction, err := h.service.GetTransactionByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed get transaction of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List of transaction campaign", http.StatusOK, "success", transaction.FormatCampaignTransactions(listTransaction))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	listTransactions, err := h.service.GetTransactionByUserID(userID)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed get transaction of user", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List of user's transactions", http.StatusOK, "success", transaction.FormatListUserTransactions(listTransactions))
	c.JSON(http.StatusOK, response)
}
