package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbHost = getEnv("POSTGRES_HOST", "postgresql")
var dbPort = getEnv("POSTGRES_PORT", "5432")
var dbUser = getEnv("POSTGRES_USER", "ui_test")
var dbPassword = getEnv("POSTGRES_PWD", "uiPassword5678")

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Init() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=ui_test  sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	retries := 5
	for err != nil {
		log.Printf("Failed to connect to database (%d)", retries)

		if retries > 1 {
			retries--
			time.Sleep(5 * time.Second)
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			continue
		}
		log.Panicln(err)
	}

	return db
}
