package models

import (
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Expense_Name string    `json:"expense_name"`
	Amount       float64   `json:"amount"`
	Date         time.Time `json:"-"`
}

func MigrateExpense(db *gorm.DB) error {
	err := db.AutoMigrate(&Expense{})

	return err
}
