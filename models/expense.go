package models

import (
	"gorm.io/gorm"
)

type Expense struct {
	ID           int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Expense_Name string  `json:"expense_name"`
	Amount       float64 `json:"amount"`
	//Created_At time.Time `json:"created_at"`
}

func MigrateExpense(db *gorm.DB) error {
	err := db.AutoMigrate(&Expense{})

	return err
}
