-- name: CreateCoupons :execrows
INSERT INTO coupons (user_id, code, discount) VALUES (?, ?, ?);

-- name: GetCouponsByCode :many
SELECT * FROM coupons WHERE code = ? FOR UPDATE;

-- name: GetUnusedCouponByUserIDAndCode :one
SELECT * FROM coupons WHERE user_id = ? AND code = ? AND used_by IS NULL FOR UPDATE;

-- name: CreateCouponsWihtCalc :execrows
INSERT INTO coupons (user_id, code, discount) VALUES (?, CONCAT(?, '_', FLOOR(UNIX_TIMESTAMP(NOW(3))*1000)), ?);

-- name: GetCouponByUsedBy :one
SELECT * FROM coupons WHERE used_by = ? LIMIT 1;

-- name: GetCouponsByUserIDAndCode :one
SELECT * FROM coupons WHERE user_id = ? AND code = ? AND used_by IS NULL LIMIT 1;

-- name: GetOldestUnusedCouponByUserID :one
SELECT * FROM coupons WHERE user_id = ? AND used_by IS NULL ORDER BY created_at LIMIT 1;
