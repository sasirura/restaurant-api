// Package services
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sasirura/restaurant-api/internal/logger"
	"github.com/sasirura/restaurant-api/internal/models"
	"github.com/square/square-go-sdk"
	"github.com/square/square-go-sdk/client"
	"gorm.io/gorm"
)

type SquareService struct {
	db     *gorm.DB
	Logger *logger.Logger
}

func New(db *gorm.DB, log *logger.Logger) *SquareService {
	service := &SquareService{
		db:     db,
		Logger: log,
	}

	var restaurant []models.Restaurant
	err := db.Find(&restaurant).Error
	if err != nil {
		log.Error("Failed to load restaurants", "error", err)
		return service
	}

	return service
}

// CreateOrder creates a new order
func (s *SquareService) CreateOrder(ctx context.Context, restaurant models.Restaurant, client *client.Client, tableNumber string, items []models.OrderItem) (*models.Order, error) {

	// OrderRequst
	lineItems := make([]*square.OrderLineItem, len(items))
	for i, item := range items {
		quantity := fmt.Sprintf("%d", item.Quantity)
		lineItems[i] = &square.OrderLineItem{
			Name:     &item.Name,
			Quantity: quantity,
			BasePriceMoney: &square.Money{
				Amount:   square.Int64(int64(item.UnitPrice)),
				Currency: square.CurrencyUsd.Ptr(),
			},
		}
	}

	createOrderReq := &square.CreateOrderRequest{
		Order: &square.Order{
			LocationID: restaurant.LocationID,
			LineItems:  lineItems,
		},
	}

	resp, err := client.Orders.Create(ctx, createOrderReq)
	if err != nil {
		s.Logger.Error("Failed to create square order", "error", err, "restaurant_id", restaurant.ID)
		return nil, fmt.Errorf("failed to create square order: %w", err)
	}

	totalPaidAmount := *resp.Order.TotalMoney.Amount - *resp.Order.NetAmountDueMoney.Amount
	order := &models.Order{
		ID:           *resp.Order.ID,
		RestautantID: restaurant.ID,
		TableNumber:  tableNumber,
		IsClosed:     *resp.Order.State == square.OrderStateCompleted,
		Items:        items,
		OpenAt:       time.Now(),
		Totals: models.OrderTotals{
			OrderID:   *resp.Order.ID,
			Discounts: float64(*resp.Order.TotalDiscountMoney.Amount),
			Due:       float64(*resp.Order.NetAmountDueMoney.Amount),
			Paid:      float64(totalPaidAmount),
			Tips:      float64(*resp.Order.TotalTipMoney.Amount),
			Total:     float64(*resp.Order.TotalMoney.Amount),
		},
	}

	err = s.db.Create(order).Error
	if err != nil {
		s.Logger.Error("Failed to save order to database", "error", err, "order_id", order.ID)
	}

	s.Logger.Info("Order created", "order_id", order.ID, "restautant_id", restaurant.ID)
	return order, nil
}

// GetOrdersByTable retrieves orders by table number
func (s *SquareService) GetOrdersByTable(ctx context.Context, restaurant models.Restaurant, tableNumber string) ([]models.Order, error) {
	var orders []models.Order
	if err := s.db.Where(&models.Order{
		RestautantID: restaurant.ID,
		TableNumber:  tableNumber,
	}).Preload("Items").Preload("Totals").Find(&orders).Error; err != nil {
		s.Logger.Error("Failed to fetch orders by table", "error", err, "table_number", tableNumber)
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	s.Logger.Info("Fetched orders by table", "table_number", tableNumber, "count", len(orders))
	return orders, nil
}

// GetOrderByID retrieves an order by ID
func (s *SquareService) GetOrderByID(ctx context.Context, restaurant models.Restaurant, orderID string) (*models.Order, error) {
	var order models.Order
	if err := s.db.Where(&models.Order{
		ID:           orderID,
		RestautantID: restaurant.ID,
	}).Preload("Items").Preload("Totals").First(&order).Error; err != nil {
		s.Logger.Error("Failed to fetch order by ID", "error", err, "order_id", orderID)
		return nil, fmt.Errorf("order not found: %w", err)
	}

	s.Logger.Info("Fetched order by ID", "order_id", orderID)
	return &order, nil
}

// ProcessPayment processes a payment for an order
func (s *SquareService) ProcessPayment(ctx context.Context, restaurant models.Restaurant, client *client.Client, orderID string, req models.PaymentRequest) error {

	// Create payment request
	paymentReq := &square.CreatePaymentRequest{
		IdempotencyKey: req.PaymentID,
		SourceID:       "CASH",
		OrderID:        &orderID,
		AmountMoney: &square.Money{
			Amount:   square.Int64(int64(req.BillAmount)),
			Currency: square.CurrencyUsd.Ptr(),
		},
		TipMoney: &square.Money{
			Amount:   square.Int64(int64(req.TipAmount)),
			Currency: square.CurrencyUsd.Ptr(),
		},
		CashDetails: &square.CashPaymentDetails{
			BuyerSuppliedMoney: &square.Money{
				Amount:   square.Int64(int64(req.BillAmount)),
				Currency: square.CurrencyUsd.Ptr(),
			},
		},
		LocationID: &restaurant.LocationID,
	}

	resp, err := client.Payments.Create(ctx, paymentReq)
	if err != nil {
		s.Logger.Error("Failed to process payment", "error", err, "order_id", orderID)
		return fmt.Errorf("failed to process payment: %w", err)
	}

	// Update order in database
	var order models.Order
	if err := s.db.Where(&models.Order{
		ID: orderID,
	}).First(&order).Error; err != nil {
		s.Logger.Error("Order not found in database", "error", err, "order_id", orderID)
		return fmt.Errorf("order not found: %w", err)
	}

	if resp.Payment.AmountMoney != nil {
		order.Totals.Paid += float64(*resp.Payment.AmountMoney.Amount)
	}
	if resp.Payment.TipMoney != nil {
		order.Totals.Tips += float64(*resp.Payment.TipMoney.Amount)
	}
	if order.Totals.Paid >= order.Totals.Due {
		order.IsClosed = true
	}

	if err := s.db.Save(&order).Error; err != nil {
		s.Logger.Error("Failed to update order in database", "error", err, "order_id", orderID)
		return fmt.Errorf("failed to update order: %w", err)
	}

	s.Logger.Info("Payment processed", "order_id", orderID, "payment_id", req.PaymentID)
	return nil
}
