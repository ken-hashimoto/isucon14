-- name: CreateChairs :execrows
INSERT INTO chairs (id, owner_id, name, model, is_active, access_token) VALUES (?, ?, ?, ?, ?, ?)