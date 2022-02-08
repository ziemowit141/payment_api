package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Id               string `gorm:"primaryKey"`
	Amount           float32
	Currency         string
	CreditCardNumber string
	Captures         []Capture `gorm:"foreignKey:TransactionId;references:Id;constraint:OnDelete:CASCADE"`
	Refunds          []Refund  `gorm:"foreignKey:TransactionId;references:Id;constraint:OnDelete:CASCADE"`
}

func (a Transaction) String() string {
	return fmt.Sprintf("Transaction<%s %s %v %s>", a.Id, a.CreditCardNumber, a.Amount, a.Currency)
}
