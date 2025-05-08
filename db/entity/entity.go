package entity

import (
	"time"
)

type Account struct {
	AccountID   uint      `gorm:"primaryKey" json:"account_id"`
	Name        string    `gorm:"not null" json:"name"`
	AccountType string    `gorm:"column:account_type;not null" json:"account_type"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

type Category struct {
	CategoryID   uint      `gorm:"primaryKey" json:"category_id"`
	Name         string    `gorm:"unique;not null" json:"name"`
	CategoryType string    `gorm:"column:category_type;not null" json:"category_type"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

type Transaction struct {
	TransactionID   uint      `gorm:"primaryKey" json:"transaction_id"`
	AccountID       uint      `gorm:"not null" json:"account_id"`
	CategoryID      uint      `gorm:"not null" json:"category_id"`
	Amount          float64   `gorm:"type:numeric(10,2);not null" json:"amount"`
	TransactionDate time.Time `gorm:"type:date;not null" json:"transaction_date"`
	Description     *string   `json:"description,omitempty"` // nullable
	CreatedAt       time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time `gorm:"not null;default:now()" json:"updated_at"`

	Account  *Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
