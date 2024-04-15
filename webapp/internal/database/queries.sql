

-- name: CreateUser :one
INSERT INTO users
(name, username, encrypted_password, age, dob, photo_path)
VALUES
(@name::text, @username::text, @encrypted_password::text, @age::integer, @dob::date, @photo_path::text)
RETURNING *;

-- name: SearchUsers :many
SELECT *
  FROM users t
  WHERE t.name ILIKE '%' || @query::text || '%'
  ORDER BY t.name ASC
  LIMIT @page_limit::int
  OFFSET @page_offset::int;

-- name: CountSearchedUsers :many
SELECT COUNT(id)
  FROM users t
  WHERE t.name ILIKE '%' || @query::text || '%';

-- name: FetchUserByID :one
SELECT *
  FROM users t
  WHERE id = @id::uuid
  LIMIT 1;

-- name: FetchUsersByIDs :many
SELECT *
  FROM users t
  WHERE id = ANY(@ids::uuid[]);

-- name: DeleteUser :exec
DELETE FROM users t
  WHERE id = @id::uuid;

-- name: UpdateUserName :one
UPDATE users t
SET name = @name::text
WHERE id = @id::uuid
RETURNING *;

-- name: UpdateUserAge :one
UPDATE users t
SET age = @age::integer
WHERE id = @id::uuid
RETURNING *;

-- name: UpdateUserDob :one
UPDATE users t
SET dob = @dob::date
WHERE id = @id::uuid
RETURNING *;

-- name: UpdateUserPhotoPath :one
UPDATE users t
SET photo_path = @photo_path::text
WHERE id = @id::uuid
RETURNING *;
