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

func (s *Session) AuthorizeWithCardDetails(ar *io_structures.AuthorizationRequest) string {
	creditCard, err := s.creditCardRepo.SelectCard(ar.CreditCardNumber)
	if err != nil {
		return "WRONG CARD NUMBER"
	}

	if ar.CreditCardCVV != creditCard.CVV {
		return "WRONG CVV"
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
	for indx_1, transaction := range s.creditCard.Transactions {
		log.Printf("Transaction nr %d, value: %v", indx_1, transaction.Amount)
		netBalance -= transaction.Amount
		for indx_2, capture := range transaction.Captures {
			log.Printf("Capture nr %d, value: %v", indx_2, capture.Amount)
			netBalance += capture.Amount
		}

		for indx_3, refund := range transaction.Refunds {
			log.Printf("Refund nr %d, value: %v", indx_3, refund.Amount)
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
