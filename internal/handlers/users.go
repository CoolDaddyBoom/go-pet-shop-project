package handlers

import (
	"go-pet-shop/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Users interface {
	GetAllUsers() ([]models.User, error)
	CreateUser(user models.User) error
	GetUserByEmail(email string) (models.User, error)
}

func GetAllUsers(log *slog.Logger, users Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.users.GetAllUsers"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		users, err := users.GetAllUsers()

		if err != nil {
			log.Error("failed to get users", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Retrieved users successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, users)
	}
}

func CreateUser(log *slog.Logger, users Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.users.CreateUser"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Creating new user", slog.String("url", r.URL.String()))

		var user models.User
		if err := render.DecodeJSON(r.Body, &user); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := users.CreateUser(user); err != nil {
			log.Error("failed to create user", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("User created successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]string{"status": "User created successfully"})
	}
}

func GetUserByEmail(log *slog.Logger, users Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.users.GetUserByEmail"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		email := chi.URLParam(r, "id")
		if email == "" {
			log.Error("empty id")
			http.Error(w, "invalid request: id is required", http.StatusBadRequest)
			return
		}

		user, err := users.GetUserByEmail(email)

		if err != nil {
			log.Error("failed to get user", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("User is retrieved successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, user)
	}
}
