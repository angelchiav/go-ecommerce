package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/angelchiav/go-ecommerce/internal/httpx"
	"github.com/angelchiav/go-ecommerce/internal/service"
)

type Cart struct {
	cart *service.CartService
}

func NewCart(cart *service.CartService) *Cart { return &Cart{cart: cart} }

func userIDFromRequest(r *http.Request) int64 {
	return httpx.MustUserID(r)
}

func (h *Cart) Get(w http.ResponseWriter, r *http.Request) {
	cv, err := h.cart.Get(r.Context(), userIDFromRequest(r))
	if err != nil {
		log.Printf("GET /v1/cart error: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "server_error")
		return
	}
	httpx.JSON(w, http.StatusOK, cv)
}

type addItemReq struct {
	ProductID int64 `json:"product_id"`
	Qty       int32 `json:"qty"`
}

func (h *Cart) AddItem(w http.ResponseWriter, r *http.Request) {
	var req addItemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.ProductID <= 0 || req.Qty <= 0 {
		httpx.Error(w, http.StatusBadRequest, "product_id_and_qty_required")
		return
	}

	err := h.cart.AddItem(r.Context(), userIDFromRequest(r), req.ProductID, req.Qty)
	if err == service.ErrQtyInvalid {
		httpx.Error(w, http.StatusBadRequest, "qty_invalid")
		return
	}
	if err != nil {
		log.Printf("POST /v1/cart/items error: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "server_error")
		return
	}

	httpx.JSON(w, http.StatusCreated, map[string]any{"status": "ok"})
}

type updateQtyReq struct {
	Qty int32 `json:"qty"`
}

func (h *Cart) UpdateItemQty(w http.ResponseWriter, r *http.Request) {
	idStr := httpx.Param(r, "id")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_item_id")
		return
	}

	var req updateQtyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_json")
		return
	}

	if req.Qty <= 0 {
		httpx.Error(w, http.StatusBadRequest, "qty_invalid")
		return
	}

	err = h.cart.UpdateItemQty(r.Context(), userIDFromRequest(r), itemID, req.Qty)
	if err == service.ErrQtyInvalid {
		httpx.Error(w, http.StatusBadRequest, "qty_invalid")
		return
	}

	if err != nil {
		log.Printf("PATCH /v1/cart/items/{id} error: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "server_error")
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (h *Cart) DeleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := httpx.Param(r, "id")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || itemID <= 0 {
		httpx.Error(w, http.StatusBadRequest, "invalid_item_id")
		return
	}

	if err := h.cart.DeleteItem(r.Context(), userIDFromRequest(r), itemID); err != nil {
		log.Printf("DELETE /v1/cart/items/{id} error: %v", err)
		httpx.Error(w, http.StatusInternalServerError, "server_error")
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"status": "ok"})
}
