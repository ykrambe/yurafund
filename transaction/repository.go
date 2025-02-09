package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	// Save(Transaction) (Transaction, error)
	// FindByCampaignID(int) ([]Transaction, error)
	// FindById(int) (Transaction, error)
	// Update(Transaction) (Transaction, error)

	GetTransactionByCampaignID(campaignID int) ([]Transaction, error)
	// GetTransactionByUserID(userID int) ([]Transaction, error)
	// GetTransactions() ([]Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetTransactionByCampaignID(campaignID int) ([]Transaction, error) {
	var transaction []Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("created_at desc").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

// func (r *repository) GetTransactionByUserID(userID int) ([]Transaction, error) {
// 	var transactions []Transaction

// 	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&transactions).Error
// 	if err != nil {
// 		return transactions, err
// 	}

// 	return transactions, nil
// }

// func (r *repository) GetTransactions() ([]Transaction, error) {
// 	var transactions []Transaction

// 	err := r.db.Order("created_at desc").Find(&transactions).Error
// 	if err != nil {
// 		return transactions, err
// 	}

// 	return transactions, nil
// }
