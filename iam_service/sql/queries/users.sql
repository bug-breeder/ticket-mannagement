-- name: CreateUser :one
INSERT INTO users (user_id, full_name, username, gender, birth_date, password, role)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6) RETURNING user_id;

-- name: GetUserByUsername :one
SELECT user_id, full_name, username, gender, birth_date, password, role
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT user_id, full_name, username, gender, birth_date, role
FROM users
WHERE user_id = $1
LIMIT 1;
