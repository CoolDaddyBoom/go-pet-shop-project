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

type Order_items interface {
	AddOrderItem(orderItem models.OrderItem) error
	GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error)
}

func AddOrderItem(log *slog.Logger, items Order_items) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.order_items.AddOrderItem"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var item models.OrderItem
		if err := render.DecodeJSON(r.Body, &item); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := items.AddOrderItem(item); err != nil {
			log.Error("failed to add order item", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, item)
	}
}

func GetOrderItemsByOrderID(log *slog.Logger, items Order_items) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.order_items.GetOrderItemsByOrderID"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		orderIDStr := chi.URLParam(r, "id")
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			log.Error("invalid order id", slog.String("id", orderIDStr), slog.Any("error", err))
			http.Error(w, "invalid order id", http.StatusBadRequest)
			return
		}

		orderItems, err := items.GetOrderItemsByOrderID(orderID)
		if err != nil {
			log.Error("failed to get order items", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, orderItems)
	}
}
