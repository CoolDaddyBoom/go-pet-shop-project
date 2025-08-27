package handlers

import (
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Orders interface {
	GetOrderByID(id int) (models.Order, error)
	CreateOrder(order models.Order) (int, error)
	GetOrdersByUserEmail(email string) ([]models.Order, error)
}

func GetOrderByID(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.GetOrderByID"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("invalid id format", slog.String("id", idStr), slog.Any("error", err))
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		order, err := orders.GetOrderByID(id)
		if err != nil {
			log.Error("failed to get order", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, order)
	}
}

func CreateOrder(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.CreateOrder"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var order models.Order
		if err := render.DecodeJSON(r.Body, &order); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := orders.CreateOrder(order)
		if err != nil {
			log.Error("failed to create order", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Order created successfully", slog.Int("id", id))
		render.JSON(w, r, map[string]int{"id": id})
	}
}

func GetOrdersByUserEmail(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.GetOrdersByUserEmail"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		email := chi.URLParam(r, "email")
		if email == "" {
			log.Error("empty email")
			http.Error(w, "email is required", http.StatusBadRequest)
			return
		}

		orderList, err := orders.GetOrdersByUserEmail(email)
		if err != nil {
			log.Error("failed to get orders", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, orderList)
	}
}
