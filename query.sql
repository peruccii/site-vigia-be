-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = $1;

-- name: RegisterUser :exec
INSERT INTO users (name, email, password_hash, email_verified_at)
VALUES ($1, $2, $3, $4);

