-- name: CreateOrder :one
INSERT INTO orders (
    id,
    coupon_code,
    created_at
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: AddProductToOrder :one
INSERT INTO order_product (
    id,
    order_id,
    product_id
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;
