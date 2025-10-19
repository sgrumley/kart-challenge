-- Get Coupon by ID
-- name: GetCouponByID :one
SELECT id 
FROM coupons
WHERE id = $1;

