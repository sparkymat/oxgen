-- name: FetchUserByUsername :one
SELECT * FROM users
WHERE username = @username::text LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  name, username, encrypted_password
) VALUES (
  @name::text, @username::text, @encrypted_password::text
) RETURNING *;
