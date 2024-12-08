-- name: GetRidesByUserIDOrderedByCreatedAtDesc :many
SELECT * FROM rides WHERE user_id = ? ORDER BY created_at DESC;

-- name: GetLatestRideStatusByRideID :one
SELECT status FROM ride_statuses WHERE ride_id = ? ORDER BY created_at DESC LIMIT 1;

-- name: GetRidesByUserID :many
SELECT * FROM rides WHERE user_id = ?;

-- name: CreateRide :exec
INSERT INTO rides (id, user_id, pickup_latitude, pickup_longitude, destination_latitude, destination_longitude) VALUES (?, ?, ?, ?, ?, ?);

-- name: CreateRideStatus :exec
INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?);

-- name: CountRidesByUserID :one
SELECT COUNT(*) FROM rides WHERE user_id = ?;
