-- name: GetRidesByUserID :many
SELECT * FROM rides WHERE user_id = ? ORDER BY created_at DESC;
