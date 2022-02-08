package models

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"gorm.io/gorm"
)

type CreditCard struct {
	gorm.Model
	Number       string `gorm:"primaryKey;unique"`
	Expiry       time.Time
	CVV          string
	Balance      float32
	BaseCurrency string
	Transactions []Transaction `gorm:"foreignKey:CreditCardNumber;references:Number;constraint:OnDelete:CASCADE"`
}

func (a *CreditCard) String() string {
	return fmt.Sprintf("CreditCard<%s %s %v %s>", a.Number, a.CVV, a.Balance, a.BaseCurrency)
}

func (a *CreditCard) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(a)
}
