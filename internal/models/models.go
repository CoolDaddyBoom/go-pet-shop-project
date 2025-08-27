package models

import "time"

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int // количество на складе
}

type User struct {
	ID    int
	Name  string
	Email string
}

type Order struct {
	ID         int
	UserEmail  string
	TotalPrice float64
	CreatedAt  time.Time
}

type OrderItem struct {
	ID        int
	OrderID   int
	ProductID int
	Quantity  int
}

type Transactions struct {
	ID        int
	OrderID   int
	Amount    float64
	Status    string
	CreatedAt time.Time
}
