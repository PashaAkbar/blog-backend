package connection

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() (DB *gorm.DB, err error) {
	dsn := "host=localhost user=postgres password=postgres dbname=blog-go port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
