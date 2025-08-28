// package handlers

// import (
// 	"go-pet-shop/internal/domain"
// 	"net/http"

// 	"github.com/go-chi/render"
// )

// type ProductsHandler struct {
// 	repo domain.ProductRepository
// }

// func NewProductsHandler(repo domain.ProductRepository) *ProductsHandler {
// 	return &ProductsHandler{repo: repo}
// }

// func (h *ProductsHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
// 	products, err := h.repo.GetAll()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	render.JSON(w, r, products)
// }

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

type Products interface {
	GetAllProducts() ([]models.Product, error)
	CreateProduct(product models.Product) error
	DeleteProduct(id int) error
	UpdateProduct(product models.Product) error
	GetProductByID(id int) (models.Product, error)
	GetPopularProducts() ([]models.PopularProduct, error)
}

func GetAllProducts(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.GetAllProducts"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		products, err := products.GetAllProducts()

		if err != nil {
			log.Error("failed to get products", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Retrieved products successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, products)
	}
}

func CreateProduct(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.CreateProduct"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Creating new product", slog.String("url", r.URL.String()))

		var product models.Product
		if err := render.DecodeJSON(r.Body, &product); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := products.CreateProduct(product); err != nil {
			log.Error("failed to create product", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Product created successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]string{"status": "Product created successfully"})
	}
}

func DeleteProduct(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.DeleteProduct"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Deleting product", slog.String("url", r.URL.String()))

		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			log.Error("empty id")
			http.Error(w, "invalid request: id is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("invalid id format", slog.String("id", idStr), slog.Any("error", err))
			http.Error(w, "invalid request: id must be a number", http.StatusBadRequest)
			return
		}

		if err := products.DeleteProduct(id); err != nil {
			log.Error("failed to delete product", slog.Int("id", id), slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Deleted product successfully", slog.Int("id", id), slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]string{"status": "Product deleted successfully"})
	}
}

func UpdateProduct(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.UpdateProduct"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Updating product", slog.String("url", r.URL.String()))

		var product models.Product
		if err := render.DecodeJSON(r.Body, &product); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := products.UpdateProduct(product); err != nil {
			log.Error("failed to update product", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Product updated successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]string{"status": "Product updated successfully"})
	}
}

func GetProductByID(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.GetProductByID"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			log.Error("empty id")
			http.Error(w, "invalid request: id is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("invalid id format", slog.String("id", idStr), slog.Any("error", err))
			http.Error(w, "invalid request: id must be a number", http.StatusBadRequest)
			return
		}

		product, err := products.GetProductByID(id)

		if err != nil {
			log.Error("failed to get a product", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("A product is retrieved successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, product)
	}
}

func GetPopularProducts(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.GetPopularProducts"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		products, err := products.GetPopularProducts()
		if err != nil {
			log.Error("failed to get popular products", slog.Any("error", err))
			http.Error(w, "failed to get popular products", http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, products)
	}
}
