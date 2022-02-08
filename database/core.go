package database

import (
	"fmt"

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

	return db
}
