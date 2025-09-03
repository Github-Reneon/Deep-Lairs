package main

import (
	"deep_lairs/internal/protocol"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDBO() {
	var err error
	db, err = gorm.Open(sqlite.Open(protocol.DBO_SQLITE), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

func createTables() {
	// Ensure db is initialized
	if db == nil {
		panic("database connection is not initialized")
	}
}
