-- name: UpdatePaymentGatewayURLValue :execrows
UPDATE settings SET value = ? WHERE name = 'payment_gateway_url';

-- name: CreatePaymentToken :execrows
INSERT INTO payment_tokens (user_id, token) VALUES (?, ?);
