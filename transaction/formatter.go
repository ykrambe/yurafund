package transaction

type CampaignTransactionFormatter struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{} // untuk membuat data formatter dengan data default (mengikuti data stuctnya)
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt.String()
	return formatter
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionFormatter {
	if len(transactions) == 0 {
		return []CampaignTransactionFormatter{}
	}

	listOfFormatterTransaction := []CampaignTransactionFormatter{}

	for _, transactions := range transactions {
		formatterTransaction := FormatCampaignTransaction(transactions)
		listOfFormatterTransaction = append(listOfFormatterTransaction, formatterTransaction)
	}

	return listOfFormatterTransaction
}
