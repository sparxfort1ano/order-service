package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sparxfort1ano/order-service/internal/cache"
	"github.com/sparxfort1ano/order-service/internal/repository"
)

type OrderHandler struct {
	cache *cache.Cache
	repo  *repository.PostgresRepository
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Any URL
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Check if id in query
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Cache order check
	order, ok := h.cache.Get(id)
	if ok {
		log.Printf("Order %s found in cache", id)
		renderJSON(w, order)
		return
	}

	// DB order check
	order, err := h.repo.GetOrderById(r.Context(), id)
	if err != nil {
		log.Printf("Order %s not found in DB: %v", id, err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Cache updating
	h.cache.Set(order)
	log.Printf("Order %s loaded from DB and cached", id)
	renderJSON(w, order)
}

// Init handler
func NewOrderHandler(c *cache.Cache, r *repository.PostgresRepository) *OrderHandler {
	return &OrderHandler{
		cache: c,
		repo:  r,
	}
}

// Auxiliary func
func renderJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
