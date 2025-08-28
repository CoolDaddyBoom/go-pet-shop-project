package handlers

import (
	"encoding/json"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Place_Orders interface {
	PlaceOrder(userEmail string, items []models.OrderItem) (int, error)
}

func PlaceOrder(log *slog.Logger, pl Place_Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.place_orders.PlaceOrder"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req struct {
			UserEmail string
			Items     []models.OrderItem
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("invalid request body", slog.Any("error", err))
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		orderID, err := pl.PlaceOrder(req.UserEmail, req.Items)
		if err != nil {
			log.Error("failed to place order", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]any{
			"status":  "success",
			"orderID": orderID,
		}

		render.JSON(w, r, resp)
	}
}
