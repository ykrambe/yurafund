package transaction

import "yurafund/user"

type GetCampaignTransactionInput struct {
	ID   int `uri:"id" binding:"required"` // mengambil parameter id pada path api
	User user.User
}
