package seed

import (
	"github.com/ziemowit141/payment_api/database/models"
	"gorm.io/gorm"
)

func LoadTestCreditCards(db *gorm.DB) {
	creditCard1 := &models.CreditCard{
		Number:       "5105105105105100",
		CVV:          "111",
		Balance:      10000.0,
		BaseCurrency: "PLN",
	}
	result := db.Create(creditCard1)

	if result.Error != nil {
		panic(result.Error)
	}

	creditCard2 := &models.CreditCard{
		Number:       "5105105105105101",
		CVV:          "111",
		Balance:      10000.0,
		BaseCurrency: "PLN",
	}
	result = db.Create(creditCard2)

	if result.Error != nil {
		panic(result.Error)
	}

	creditCard3 := &models.CreditCard{
		Number:       "5105105105105102",
		CVV:          "111",
		Balance:      10000.0,
		BaseCurrency: "PLN",
	}
	result = db.Create(creditCard3)

	if result.Error != nil {
		panic(result.Error)
	}

	creditCard4 := &models.CreditCard{
		Number:       "5105105105105103",
		CVV:          "111",
		Balance:      10000.0,
		BaseCurrency: "PLN",
	}
	result = db.Create(creditCard4)

	if result.Error != nil {
		panic(result.Error)
	}

	creditCard5 := &models.CreditCard{
		Number:       "5105105105105120",
		CVV:          "111",
		Balance:      10000.0,
		BaseCurrency: "PLN",
	}
	result = db.Create(creditCard5)

	if result.Error != nil {
		panic(result.Error)
	}
}

func LoadTestTransaciton(db *gorm.DB) {
	transaction := &models.Transaction{
		Id:               "35c64636-70bc-42b5-593e-d051e29c02bc",
		Amount:           5000.0,
		Currency:         "PLN",
		CreditCardNumber: "5105105105105120",
	}
	result := db.Create(transaction)

	if result.Error != nil {
		panic(result.Error)
	}
}
