package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/models"
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
