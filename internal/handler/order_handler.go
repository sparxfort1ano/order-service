package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sparxfort1ano/order-service/internal/service"
)

// HTTP-обработчик запросов
type Handler struct{ svc *service.Service }

// Конструктор-обработчик с привязанным сервисом
func New(s *service.Service) *Handler { return &Handler{svc: s} }

// Регистрируем эндпоинты для HTTP-сервера
func (h *Handler) Routes(mux *http.ServeMux) {
	// отдаем index.html
	mux.Handle("/", http.FileServer(http.Dir("web/static")))

	// эндпоинт получения заказа по id -> GET /order/{id}
	mux.HandleFunc("GET /order/", h.getOrder)

	// простая проверка на то, что сервис живой -> GET /ping
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("pong")) })
}

// Отдаем заказ в JSON по id
func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	// достаем id из URL
	id := strings.TrimPrefix(r.URL.Path, "/order/")
	if id == "" || strings.ContainsRune(id, '/') {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}

	// пробуем получить заказ через сервис
	o, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// отправляем JSON
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(o)
}
