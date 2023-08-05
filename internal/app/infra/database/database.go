package database

import (
	"log"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDatabase() *gorm.DB {
	dsn := "host=db user=postgres password=postgres dbname=rinha port=5432"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.AutoMigrate(&peopledb.PersonModel{})

	return db
}
