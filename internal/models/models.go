// Package models
package models

import (
	"time"

	"gorm.io/gorm"
)

type Restaurant struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	Name        string
	SqaureToken string `gorm:"unique"`
	LocationID  string
}

type Order struct {
	gorm.Model
	ID           string `gorm:"primaryKey"`
	RestautantID uint
	TableNumber  string
	IsClosed     bool
	Items        []OrderItem `gorm:"foreignKey:OrderID"`
	Totals       OrderTotals `gorm:"foreignKey:OrderID"`
	OpenAt       time.Time
}

type OrderItem struct {
	gorm.Model
	OrderID   string
	Name      string
	Comment   string
	UnitPrice float64
	Quantity  int
	Discounts []Discount `gorm:"foreignKey:OrderItemID"`
	Modifiers []Modifier `gorm:"foreignKey:OrderItemID"`
	Amount    float64
}

type Discount struct {
	gorm.Model
	OrderItemID  uint
	Name         string
	IsPercentage bool
	Value        float64
	Amount       float64
}

type Modifier struct {
	gorm.Model
	OrderItemID uint
	Name        string
	UnitPrice   float64
	Quantity    int
	Amount      float64
}

type OrderTotals struct {
	gorm.Model
	OrderID       string
	Discounts     float64
	Due           float64
	Tax           float64
	ServiceCharge float64
	Paid          float64
	Tips          float64
	Total         float64
}

type PaymentRequest struct {
	BillAmount float64 `json:"billAmount"`
	TipAmount  float64 `json:"tipAmount"`
	PaymentID  string  `json:"paymentId"`
}
