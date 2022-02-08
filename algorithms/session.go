package algorithms

import (
	"log"

	"github.com/ziemowit141/payment_api/database/models"
	"github.com/ziemowit141/payment_api/database/repositories"
	"github.com/ziemowit141/payment_api/handlers/io_structures"
	"gorm.io/gorm"
)

type Session struct {
	creditCardRepo  *repositories.CreditCardRepository
	transactionRepo *repositories.TrasnsactionRepository
	captureRepo     *repositories.CaptureRepository
	refundRepo      *repositories.RefundRepository

	creditCard  *models.CreditCard
	transaction *models.Transaction
}

func NewSession(db *gorm.DB) *Session {
	return &Session{repositories.NewCreditCardRepository(db),
		repositories.NewTrasnsactionRepository(db),
		repositories.NewCaptureRepository(db),
		repositories.NewRefundRepository(db),
		nil,
		nil}
}
