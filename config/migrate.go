package config

import (
	"log"
	"task-5-pbi-btpns-daffa_satria/models"
)

func Migrate() {
	db.AutoMigrate(&models.User{}, &models.Photo{})
	log.Println("Database Migration Completed!")
}