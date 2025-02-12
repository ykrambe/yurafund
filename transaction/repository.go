package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetTransactionByCampaignID(campaignID int) ([]Transaction, error)
	GetTransactionByUserID(userID int) ([]Transaction, error)
	GetTransactions() ([]Transaction, error)
	GetTransactionByID(ID int) (Transaction, error)
	SaveTransaction(transaction Transaction) (Transaction, error)
	UpdateTransaction(transaction Transaction) (Transaction, error)
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

func (r *repository) GetTransactionByUserID(userID int) ([]Transaction, error) {
	var transactions []Transaction
	// preload 2 kali untuk mengambil data campaign lalu campaign imagesnya
	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) GetTransactions() ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("Campaign").Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) SaveTransaction(transaction Transaction) (Transaction, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (r *repository) UpdateTransaction(transaction Transaction) (Transaction, error) {
	err := r.db.Save(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (r *repository) GetTransactionByID(ID int) (Transaction, error) {
	var transaction Transaction
	// preload 2 kali untuk mengambil data campaign lalu campaign imagesnya
	err := r.db.Where("id = ?", ID).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
