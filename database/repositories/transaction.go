package repositories

import (
	"github.com/ziemowit141/payment_api/database/models"
	"github.com/ziemowit141/payment_api/util"
	"gorm.io/gorm"
)

type TrasnsactionRepository struct {
	db *gorm.DB
}

func NewTrasnsactionRepository(db *gorm.DB) *TrasnsactionRepository {
	return &TrasnsactionRepository{db}
}

func (tr *TrasnsactionRepository) NewTransaction(amount float32, creditCard *models.CreditCard) *models.Transaction {
	transaction := &models.Transaction{
		Id:               util.GenerateUniqeId(),
		Amount:           amount,
		Currency:         creditCard.BaseCurrency,
		CreditCardNumber: creditCard.Number,
	}

	result := tr.db.Create(transaction)
	if result.Error != nil {
		panic(result.Error)
	}

	err := tr.db.Model(creditCard).Association("Transactions").Append(transaction)
	if err != nil {
		panic(err)
	}

	return transaction
}

func (tr *TrasnsactionRepository) SelectTransaction(uid string) *models.Transaction {
	transaction := new(models.Transaction)
	result := tr.db.Preload("Captures").Preload("Refunds").First(transaction, "id = ?", uid)

	if result.Error != nil {
		return nil
	}

	return transaction
}

func (tr *TrasnsactionRepository) DeleteTransaction(transaction *models.Transaction) {
	result := tr.db.Delete(transaction, "id = ?", transaction.Id)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (tr *TrasnsactionRepository) RefundCount(transaction *models.Transaction) int64 {
	return tr.db.Model(transaction).Association("Refunds").Count()
}
