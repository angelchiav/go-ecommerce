package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/angelchiav/go-ecommerce/internal/sqlc"
)

var (
	ErrQtyInvalid   = errors.New("qty_invalid")
	ErrItemNotFound = errors.New("item_not_found")
)

type CartItem struct {
	ID             int64  `json:"id"`
	ProductID      int64  `json:"product_id"`
	Name           string `json:"name"`
	Qty            int32  `json:"qty"`
	PriceCents     int32  `json:"price_cents"`
	LineTotalCents int32  `json:"line_total_cents"`
}

type CartView struct {
	CartID     int64      `json:"cart_id"`
	Items      []CartItem `json:"items"`
	TotalCents int32      `json:"total_cents"`
}

type CartService struct {
	q  *sqlc.Queries
	db *sql.DB
}

func NewCartService(db *sql.DB, q *sqlc.Queries) *CartService {
	return &CartService{db: db, q: q}
}

func (s *CartService) Get(ctx context.Context, userID int64) (*CartView, error) {
	cartID, err := s.q.GetOrCreateActiveCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	rows, err := s.q.ListCartItems(ctx, cartID)
	if err != nil {
		return nil, err
	}

	items := make([]CartItem, 0, len(rows))
	var total int32

	for _, r := range rows {
		items = append(items, CartItem{
			ID:             r.ID,
			ProductID:      r.ProductID,
			Name:           r.Name,
			Qty:            r.Qty,
			PriceCents:     int32(r.PriceCents),
			LineTotalCents: r.LineTotalCents,
		})
		total += r.LineTotalCents
	}
	return &CartView{
		CartID:     cartID,
		Items:      items,
		TotalCents: total,
	}, nil
}

func (s *CartService) AddItem(ctx context.Context, userID, productID int64, qty int32) error {
	if qty <= 0 {
		return ErrQtyInvalid
	}
	cartID, err := s.q.GetActiveCartID(ctx, userID)
	if err != nil {
		return err
	}
	_, err = s.q.UpsertCartItem(ctx, sqlc.UpsertCartItemParams{
		CartID:    cartID,
		ProductID: productID,
		Qty:       qty,
	})
	return err
}

func (s *CartService) UpdateItemQty(ctx context.Context, userID, itemID int64, qty int32) error {
	if qty <= 0 {
		return ErrQtyInvalid
	}
	cartID, err := s.q.GetActiveCartID(ctx, userID)
	if err != nil {
		return err
	}
	if err := s.q.UpdateCartItemQtyInCart(ctx, sqlc.UpdateCartItemQtyInCartParams{
		ID:     itemID,
		CartID: cartID,
		Qty:    qty,
	}); err != nil {
		return err
	}
	return nil
}

func (s *CartService) DeleteItem(ctx context.Context, userID, itemID int64) error {
	cartID, err := s.q.GetActiveCartID(ctx, userID)
	if err != nil {
		return err
	}
	if err := s.q.DeleteCartItemInCart(ctx, sqlc.DeleteCartItemInCartParams{
		ID:     itemID,
		CartID: cartID,
	}); err != nil {
		return err
	}
	return nil
}
