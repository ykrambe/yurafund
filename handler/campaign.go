package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"yurafund/campaign"
	"yurafund/helper"
	"yurafund/user"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		fmt.Println(err)
		responseError := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		responseError := helper.APIResponse("Error to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}
	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		responseError := helper.APIResponse("Error to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}

	response := helper.APIResponse("campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		responseError := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errors)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)

	if err != nil {
		responseError := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}

	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}
