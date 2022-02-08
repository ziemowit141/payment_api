package repositories

import (
	"github.com/ziemowit141/payment_api/database/models"
	"gorm.io/gorm"
)

type CreditCardRepository struct {
	db *gorm.DB
}

func NewCreditCardRepository(db *gorm.DB) *CreditCardRepository {
	return &CreditCardRepository{db}
}

func (ar *CreditCardRepository) SelectCard(ccn string) (*models.CreditCard, error) {
	acc := new(models.CreditCard)
	result := ar.db.Preload("Transactions").
		Preload("Transactions.Captures").
		Preload("Transactions.Refunds").
		First(acc, "number = ?", ccn)

	return acc, result.Error
}

func (ar *CreditCardRepository) DeleteCard(creditCard *models.CreditCard) {
	result := ar.db.Delete(creditCard, "number = ?", creditCard.Number)

	if result.Error != nil {
		panic(result.Error)
	}
}

func (ar *CreditCardRepository) UpdateCard(creditCard *models.CreditCard) {
	result := ar.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(creditCard)

	if result.Error != nil {
		panic(result.Error)
	}
}
