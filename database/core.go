package database

import (
	"fmt"

	"github.com/ziemowit141/payment_api/database/models"
	"github.com/ziemowit141/payment_api/database/seed"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase(name string) *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=localhost "+
			"user=postgres "+
			"password=mysecretpassword "+
			"dbname=%s "+
			"port=5432 "+
			"sslmode=disable ", name),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	createSchema(db)

	return db
}

func createSchema(db *gorm.DB) {
	err := db.AutoMigrate(&models.CreditCard{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Transaction{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Capture{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Refund{})
	if err != nil {
		panic(err)
	}
}

func DropSchema(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.Capture{})
	if err != nil {
		panic(err)
	}

	err = db.Migrator().DropTable(&models.Refund{})
	if err != nil {
		panic(err)
	}

	err = db.Migrator().DropTable(&models.Transaction{})
	if err != nil {
		panic(err)
	}

	err = db.Migrator().DropTable(&models.CreditCard{})
	if err != nil {
		panic(err)
	}
}

func NewTestDb() *gorm.DB {
	db := SetupDatabase("myapp")
	seed.LoadTestCreditCards(db)
	seed.LoadTestTransaciton(db)

	return db
}
