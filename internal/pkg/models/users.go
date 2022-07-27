package models

import "time"

// swagger:models Users
// @Description Full user data
type Users struct {
	// User account
	Acct string `json:"account" gorm:"primaryKey; column:acct"`
	// User's password, hashed
	Password string `json:"password" gorm:"column:pwd"`
	// User's full name
	FullName string `json:"fullName" gorm:"column:fullname"`
	// The time when the account was created
	CreatedAt time.Time `json:"createdAt"`
	// The time when the account was last updated
	UpdatedAt time.Time `json:"updatedAt"`
}

// swagger:models UsersList
// @Description Partial user data for list user API
type UsersList struct {
	// User account
	Acct string `json:"account" gorm:"primaryKey; column:acct"`
	// User's full name
	FullName string `json:"fullName" gorm:"column:fullname"`
}
