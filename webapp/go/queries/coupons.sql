-- name: CreateCoupons :execrows
INSERT INTO coupons (user_id, code, discount) VALUES (?, ?, ?);

-- name: GetCouponsByCode :many
SELECT * FROM coupons WHERE code = ? FOR UPDATE;

-- name: CreateCouponsWihtCalc :execrows
INSERT INTO coupons (user_id, code, discount) VALUES (?, CONCAT(?, '_', FLOOR(UNIX_TIMESTAMP(NOW(3))*1000)), ?)
