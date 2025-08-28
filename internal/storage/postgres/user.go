package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
)

func (s *Storage) GetAllUsers() ([]models.User, error) {
	const fn = "storage.postgres.user.GetAllUsers"

	rows, err := s.db.Query(context.Background(), `SELECT * FROM users`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Storage) CreateUser(user models.User) error {
	const fn = "storage.postgres.user.CreateUser"

	_, err := s.db.Exec(context.Background(),
		`INSERT INTO users (name, email) VALUES ($1, $2)`,
		user.Name, user.Email)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) GetUserByEmail(email string) (models.User, error) {
	const fn = "storage.postgres.user.GetUserByEmail"

	var user models.User
	err := s.db.QueryRow(context.Background(), `SELECT id, name, email FROM users 
	WHERE email = $1`, email).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", fn, err)
	}
	return user, nil
}
func (s *Storage) GetUserOrderHistory(email string) ([]models.OrderDetail, error) {
	const query = `
        SELECT 
            o.id AS order_id,
            o.user_email,
            o.total_price,
            o.created_at,
            p.name AS product_name,
            oi.quantity,
            t.id AS transaction_id,
            t.status
        FROM orders o
        JOIN order_items oi ON o.id = oi.order_id
        JOIN products p ON oi.product_id = p.id
        JOIN transactions t ON o.id = t.order_id
        WHERE o.user_email = $1
        ORDER BY o.created_at DESC
    `

	rows, err := s.db.Query(context.Background(), query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.OrderDetail
	for rows.Next() {
		var od models.OrderDetail
		if err := rows.Scan(
			&od.OrderID,
			&od.UserEmail,
			&od.TotalPrice,
			&od.CreatedAt,
			&od.ProductName,
			&od.Quantity,
			&od.TransactionID,
			&od.Status,
		); err != nil {
			return nil, err
		}
		history = append(history, od)
	}

	return history, nil
}
