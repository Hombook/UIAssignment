package db

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbHost     = getEnv("POSTGRES_HOST", "postgresql")
	dbPort     = getEnv("POSTGRES_PORT", "5432")
	dbUser     = getEnv("POSTGRES_USER", "ui_test")
	dbPassword = getEnv("POSTGRES_PWD", "uiPassword5678")
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// swagger:db Pagination
// @Description JSON response body to hold paginated data
type Pagination struct {
	// Max number of items per page(5 <= limit <= 100)
	Limit int `json:"limit"`
	// Requested page
	Page int `json:"page"`
	// Total number of the matched item
	TotalRows int64 `json:"totalRows"`
	// Total number of pages by the given limit
	TotalPages int `json:"totalPages"`
	// Content of data
	// example: [{"account": "ac1", "fullName": "mister man"}, {"account": "ac2", "fullName": "miss lady"}]
	Rows interface{} `json:"rows"`
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

// Pagination function for GORM Scopes
func Paginate(queryStruct interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(queryStruct).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		page := pagination.Page
		if page == 0 {
			page = 1
		}

		limit := pagination.Limit
		switch {
		case limit > 100:
			limit = 100
		case limit < 5:
			limit = 5
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}
