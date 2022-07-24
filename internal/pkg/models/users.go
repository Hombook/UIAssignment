package models

import "time"

type Users struct {
	Acct      string    `json:"account" gorm:"primaryKey; column:acct"`
	Password  string    `json:"password" gorm:"column:pwd"`
	FullName  string    `json:"fullName" gorm:"column:fullname"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
