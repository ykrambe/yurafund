package handler

import (
	"net/http"
	"yurafund/transaction"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	transactionService transaction.Service
}

func NewTransactionHandler(transactionService transaction.Service) *transactionHandler {
	return &transactionHandler{transactionService}
}

func (h *transactionHandler) Index(c *gin.Context) {
	transactions, err := h.transactionService.GetTransactions()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "transaction_index.html", gin.H{"transactions": transactions})
}
