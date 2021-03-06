package algorithms

import (
	"time"

	"github.com/ziemowit141/payment_api/database/models"
	"github.com/ziemowit141/payment_api/database/repositories"
	"github.com/ziemowit141/payment_api/handlers/io_structures"
	"github.com/ziemowit141/payment_api/util"
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

func (s *Session) AuthorizeWithCardDetails(ar *io_structures.AuthorizationRequest) string {
	creditCard, err := s.creditCardRepo.SelectCard(ar.CreditCardNumber)
	if err != nil {
		return "WRONG CARD NUMBER"
	}

	if ar.CreditCardCVV != creditCard.CVV {
		return "WRONG CVV"
	}

	expiryDate := util.ParseExpiryDate(ar.Expiry)
	if expiryDate.Before(time.Time(time.Now())) ||
		expiryDate.Year() != creditCard.Expiry.Year() ||
		expiryDate.Month() != creditCard.Expiry.Month() {
		return "CARD EXPIRED"
	}

	s.creditCard = creditCard

	if ar.Amount > s.CardNetBalance() {
		return "UNSUFFICIENT FUNDS"
	}

	return "SUCCESS"
}

func (s *Session) AuthorizeWithTransactionId(uid string) string {
	transaction := s.transactionRepo.SelectTransaction(uid)

	if transaction == nil {
		return "WRONG TRANSACTION ID"
	}

	s.transaction = transaction

	//Should be no errors as transaction must be bound to a CreditCard
	creditCard, _ := s.creditCardRepo.SelectCard(transaction.CreditCardNumber)
	s.creditCard = creditCard

	return "SUCCESS"
}

func (s *Session) AddTransaction(amount float32) string {
	if s.creditCard == nil {
		panic("unauthorized")
	}
	transaction := s.transactionRepo.NewTransaction(amount, s.creditCard)
	return transaction.Id
}

func (s *Session) CancelTransaction(uid string) {
	if s.transaction == nil {
		panic("unauthorized")
	}
	s.transactionRepo.DeleteTransaction(s.transaction)
	s.transaction = nil
}

// CardNetBalance or DisplayedBalance follows this equation:
// DisplayedBalance = balance - transactions - refunds + captures
func (s *Session) CardNetBalance() float32 {
	if s.creditCard == nil {
		panic("unauthorized")
	}
	// Hack, refreshes associations - transactions, captures and refunds
	s.creditCard, _ = s.creditCardRepo.SelectCard(s.creditCard.Number)
	var netBalance float32 = s.creditCard.Balance
	for _, transaction := range s.creditCard.Transactions {
		netBalance -= transaction.Amount
		for _, capture := range transaction.Captures {
			netBalance += capture.Amount
		}
		for _, refund := range transaction.Refunds {
			netBalance -= refund.Amount
		}
	}

	return netBalance
}

func (s *Session) TransactionNetBalance() float32 {
	if s.transaction == nil {
		panic("unauthorized")
	}
	var netBalance float32 = s.transaction.Amount
	for _, order := range s.transaction.Captures {
		netBalance -= order.Amount
	}

	return netBalance
}

func (s *Session) MaxRefundValue() float32 {
	if s.transaction == nil {
		panic("unauthorized")
	}
	var totalValue float32 = 0.0
	for _, order := range s.transaction.Captures {
		totalValue += order.Amount
	}

	for _, refund := range s.transaction.Refunds {
		totalValue -= refund.Amount
	}

	return totalValue
}

func (s *Session) Refund(amount float32) (string, float32) {
	if s.transaction == nil {
		panic("unauthorized")
	}

	if amount > s.MaxRefundValue() {
		return "REFUND CANT BE LARGER THAN CAPTURED VALUE", s.CardNetBalance()
	}

	s.refundRepo.NewRefund(amount, s.transaction)
	s.creditCard.Balance += amount
	s.creditCardRepo.UpdateCard(s.creditCard)

	return "SUCCESS", s.CardNetBalance()
}

func (s *Session) Capture(amount float32) (string, float32) {
	if s.transaction == nil {
		panic("unauthorized")
	}

	if s.transactionRepo.RefundCount(s.transaction) > 0 {
		return "REFUND WAS ISSUED - CAPTURES BLOCKED", s.CardNetBalance()
	}

	if s.TransactionNetBalance() < amount {
		return "CAPTURE CANT BE LARGER THAN TRANSACTION", s.CardNetBalance()
	}

	s.captureRepo.NewCapture(amount, s.transaction)
	s.creditCard.Balance -= amount
	s.creditCardRepo.UpdateCard(s.creditCard)

	return "SUCCESS", s.CardNetBalance()
}
