package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
	"time"
)

func (s *Storage) GetOrderByID(id int) (models.Order, error) {
	const fn = "storage.postgres.order.GetOrderByID"

	var order models.Order
	err := s.db.QueryRow(context.Background(),
		`SELECT id, user_email, total_price, created_at FROM orders WHERE id = $1`, id).
		Scan(&order.ID, &order.UserEmail, &order.TotalPrice, &order.CreatedAt)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", fn, err)
	}

	return order, nil
}

func (s *Storage) CreateOrder(order models.Order) (int, error) {
	const fn = "storage.postgres.order.CreateOrder"

	var id int
	err := s.db.QueryRow(context.Background(),
		`INSERT INTO orders (user_email, total_price, created_at) VALUES ($1, $2, $3) RETURNING id`,
		order.UserEmail, order.TotalPrice, time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}
	return id, nil
}

func (s *Storage) GetOrdersByUserEmail(email string) ([]models.Order, error) {
	const fn = "storage.postgres.order.GetOrdersByUserEmail"

	rows, err := s.db.Query(context.Background(),
		`SELECT id, user_email, total_price, created_at FROM orders WHERE user_email = $1`, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserEmail, &order.TotalPrice, &order.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}
