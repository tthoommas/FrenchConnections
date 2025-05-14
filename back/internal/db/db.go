package db

import (
	"FrenchConnections/internal"
	"FrenchConnections/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB_CLIENT *gorm.DB = nil

func GetDBClient() *gorm.DB {
	if DB_CLIENT == nil {
		// Load db client once
		dbClient, err := gorm.Open(sqlite.Open(internal.DB_PATH), &gorm.Config{})
		if err != nil || dbClient == nil {
			panic(err)
		}
		DB_CLIENT = dbClient
		DB_CLIENT.AutoMigrate(&models.Game{})
		DB_CLIENT.AutoMigrate(&models.GameCategory{})
	}
	return DB_CLIENT
}
