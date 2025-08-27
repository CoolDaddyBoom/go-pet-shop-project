package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
)

func (s *Storage) AddOrderItem(orderItem models.OrderItem) error {
	const fn = "storage.postgres.order_item.AddOrderItem"

	_, err := s.db.Exec(context.Background(),
		`INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`,
		orderItem.OrderID, orderItem.ProductID, orderItem.Quantity)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Storage) GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error) {
	const fn = "storage.postgres.order_item.GetOrderItemsByOrderID"

	rows, err := s.db.Query(context.Background(),
		`SELECT id, order_id, product_id, quantity FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		items = append(items, item)
	}

	return items, nil
}
