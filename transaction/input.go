package transaction

import "yurafund/user"

type GetCampaignTransactionInput struct {
	ID   int `uri:"id" binding:"required"` // mengambil parameter id pada path api
	User user.User
}

type CreateTransactionInput struct {
	Amount     int `json:"amount" binding:"required"`
	CampaignID int `json:"campaign_id" binding:"required"`
	User       user.User
}
