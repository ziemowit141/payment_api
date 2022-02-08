package repositories

import (
	"github.com/ziemowit141/payment_api/database/models"
	"gorm.io/gorm"
)

type CaptureRepository struct {
	db *gorm.DB
}

func NewCaptureRepository(db *gorm.DB) *CaptureRepository {
	return &CaptureRepository{db}
}

func (or *CaptureRepository) NewCapture(amount float32, transaction *models.Transaction) *models.Capture {
	capture := &models.Capture{
		TransactionId: transaction.Id,
		Amount:        amount,
		Currency:      transaction.Currency,
	}

	result := or.db.Create(capture)
	if result.Error != nil {
		panic(result.Error)
	}

	return capture
}

func (or *CaptureRepository) SelectCapture(uid int32) *models.Capture {
	capture := new(models.Capture)
	result := or.db.First(capture, uid)

	if result.Error != nil {
		return nil
	}

	return capture
}

func (or *CaptureRepository) DeleteCapture(capture *models.Capture) {
	result := or.db.Delete(capture, capture.Id)
	if result.Error != nil {
		panic(result.Error)
	}
}
