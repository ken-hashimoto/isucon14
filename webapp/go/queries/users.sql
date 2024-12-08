-- name: CreateUsers :execrows
INSERT INTO users (id, username, firstname, lastname, date_of_birth, access_token, invitation_code) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByInvitationCode :one
SELECT * FROM users WHERE invitation_code = ?;
