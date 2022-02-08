package models

import "gorm.io/gorm"

type Refund struct {
	gorm.Model
	Id            int32
	TransactionId string
	Amount        float32
	Currency      string
}
