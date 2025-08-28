package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
)

func (s *Storage) PlaceOrder(userEmail string, items []models.OrderItem) (int, error) {
	const fn = "storage.postgres.place_order.PlaceOrder"

	ctx := context.Background()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: begin tx: %w", fn, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// 1. Проверяем и уменьшаем stock
	for _, item := range items {
		res, err := tx.Exec(ctx,
			`UPDATE products 
             SET stock = stock - $1 
             WHERE id = $2 AND stock >= $1`,
			item.Quantity, item.ProductID,
		)
		if err != nil {
			return 0, fmt.Errorf("%s: update stock: %w", fn, err)
		}
		rowsAffected := res.RowsAffected()
		if rowsAffected == 0 {
			return 0, fmt.Errorf("%s: not enough stock for product %d", fn, item.ProductID)
		}
	}

	// 2. Создаём заказ
	var orderID int
	err = tx.QueryRow(ctx,
		`INSERT INTO orders (user_email, total_price, created_at) 
         VALUES ($1, 0, NOW()) RETURNING id`,
		userEmail,
	).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("%s: insert order: %w", fn, err)
	}

	totalPrice := 0.0

	// 3. Добавляем позиции заказа
	for _, item := range items {
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity, price) 
             VALUES ($1, $2, $3, $4)`,
			orderID, item.ProductID, item.Quantity, item.Price,
		)
		if err != nil {
			return 0, fmt.Errorf("%s: insert order item: %w", fn, err)
		}
		totalPrice += float64(item.Quantity) * item.Price
	}

	// 4. Обновляем total_price заказа
	_, err = tx.Exec(ctx,
		`UPDATE orders SET total_price = $1 WHERE id = $2`,
		totalPrice, orderID,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: update order total: %w", fn, err)
	}

	// 5. Создаём запись транзакции
	_, err = tx.Exec(ctx,
		`INSERT INTO transactions (order_id, amount, created_at) 
         VALUES ($1, $2, NOW())`,
		orderID, totalPrice,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: insert transaction: %w", fn, err)
	}

	// 6. Commit
	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: commit: %w", fn, err)
	}

	return orderID, nil
}
