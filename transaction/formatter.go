package transaction

import "time"

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type UserTransactionFormatter struct {
	ID        int                   `json:"id"`
	Amount    int                   `json:"amount"`
	Status    string                `json:"status"`
	CreatedAt time.Time             `json:"created_at"`
	Campaign  UserCampaignFormatter `json:"campaign"`
}

type UserCampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{} // untuk membuat data formatter dengan data default (mengikuti data stuctnya)
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt
	return formatter
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionFormatter {
	if len(transactions) == 0 {
		return []CampaignTransactionFormatter{}
	}

	listOfFormatterTransaction := []CampaignTransactionFormatter{}

	for _, transaction := range transactions {
		formatterTransaction := FormatCampaignTransaction(transaction)
		listOfFormatterTransaction = append(listOfFormatterTransaction, formatterTransaction)
	}

	return listOfFormatterTransaction
}

func FormatSingleUserTransaction(transactions Transaction) UserTransactionFormatter {
	formatter := UserTransactionFormatter{}
	formatter.ID = transactions.ID
	formatter.Amount = transactions.Amount
	formatter.Status = transactions.Status
	formatter.CreatedAt = transactions.CreatedAt
	campaignFormatter := UserCampaignFormatter{}
	campaignFormatter.Name = transactions.Campaign.Name
	campaignFormatter.ImageURL = ""
	if len(transactions.Campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = transactions.Campaign.CampaignImages[0].FileName
	}
	formatter.Campaign = campaignFormatter
	return formatter
}

func FormatListUserTransactions(transactions []Transaction) []UserTransactionFormatter {
	if len(transactions) == 0 {
		return []UserTransactionFormatter{}
	}

	listOfFormatter := []UserTransactionFormatter{}

	for _, transaction := range transactions {
		formatterTransaction := FormatSingleUserTransaction(transaction)
		listOfFormatter = append(listOfFormatter, formatterTransaction)
	}

	return listOfFormatter
}
