-- name: GetOwnerByID :one
SELECT * FROM owners WHERE id = ? LIMIT 1;