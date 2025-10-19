-- Get Product by ID
-- name: GetProductByID :one
SELECT id, name, category, price, created_at
FROM products
WHERE id = $1;


-- List all Products
-- name: ListProducts :many
SELECT id, name, category, price, created_at
FROM products
ORDER BY name;
