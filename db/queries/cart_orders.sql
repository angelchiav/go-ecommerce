-- name: GetOrCreateActiveCart :one
WITH existing AS (
  SELECT c.id FROM carts c WHERE c.user_id = $1 AND c.status = 'active' LIMIT 1
),
inserted AS (
  INSERT INTO carts (user_id, status)
  SELECT $1, 'active'
  WHERE NOT EXISTS (SELECT 1 FROM existing)
  RETURNING id
)
SELECT id FROM inserted
UNION ALL
SELECT id FROM existing
LIMIT 1;

-- name: ListCartItems :many
SELECT
  ci.id,
  ci.product_id,
  ci.qty,
  p.name,
  p.price_cents,
  (p.price_cents * ci.qty)::int AS line_total_cents
FROM cart_items ci
JOIN products p ON p.id = ci.product_id
WHERE ci.cart_id = $1
ORDER BY ci.id;

-- name: UpsertCartItem :one
INSERT INTO cart_items (cart_id, product_id, qty)
VALUES ($1, $2, $3)
ON CONFLICT (cart_id, product_id)
DO UPDATE SET qty = cart_items.qty + EXCLUDED.qty, updated_at = now()
RETURNING id, cart_id, product_id, qty;

-- name: UpdateCartItemQty :exec
UPDATE cart_items
SET qty = $2, updated_at = now()
WHERE id = $1;

-- name: DeleteCartItem :exec
DELETE FROM cart_items WHERE id = $1;

-- name: LockCartItemsForCheckout :many
SELECT ci.product_id, ci.qty, p.price_cents, p.stock, p.is_active
FROM cart_items ci
JOIN products p ON p.id = ci.product_id
WHERE ci.cart_id = $1
FOR UPDATE;

-- name: DecrementProductStock :exec
UPDATE products
SET stock = stock - $2, updated_at = now()
WHERE id = $1 AND stock >= $2;

-- name: CreateOrder :one
INSERT INTO orders (user_id, status, total_cents)
VALUES ($1, 'placed', $2)
RETURNING id;

-- name: CreateOrderItem :exec
INSERT INTO order_items (order_id, product_id, unit_price_cents, qty, line_total_cents)
VALUES ($1, $2, $3, $4, $5);

-- name: MarkCartCheckedOut :exec
UPDATE carts
SET status = 'checked_out', updated_at = now()
WHERE id = $1;

-- name: ClearCartItems :exec
DELETE FROM cart_items WHERE cart_id = $1;