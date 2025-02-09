package transaction

import (
	"errors"
	"yurafund/campaign"
)

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error)
	// GetTransactionByUserID(userID int) ([]Transaction, error)
	// GetTransactions() ([]Transaction, error)
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error) {
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, errors.New("campaign not found")
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("not owner of the campaign")
	}

	transaction, err := s.repository.GetTransactionByCampaignID(input.ID)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
