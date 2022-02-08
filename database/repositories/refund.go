package repositories

import (
	"github.com/ziemowit141/payment_api/database/models"
	"gorm.io/gorm"
)

type RefundRepository struct {
	db *gorm.DB
}

func NewRefundRepository(db *gorm.DB) *RefundRepository {
	return &RefundRepository{db}
}

func (rr *RefundRepository) NewRefund(amount float32, transaction *models.Transaction) *models.Refund {
	refund := &models.Refund{
		TransactionId: transaction.Id,
		Amount:        amount,
		Currency:      transaction.Currency,
	}

	result := rr.db.Create(refund)
	if result.Error != nil {
		panic(result.Error)
	}

	err := rr.db.Model(transaction).Association("Refunds").Append(refund)
	if err != nil {
		panic(err)
	}

	return refund
}

func (rr *RefundRepository) SelectRefund(uid int32) *models.Refund {
	refund := new(models.Refund)
	result := rr.db.First(refund, uid)

	if result.Error != nil {
		return nil
	}

	return refund
}

func (rr *RefundRepository) DeleteRefund(refund *models.Refund) {
	result := rr.db.Delete(refund, refund.Id)
	if result.Error != nil {
		panic(result.Error)
	}
}
