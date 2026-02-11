-- name: CreateUser :one
INSERT INTO users (id,created_at,updated_at,name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUserByName :one

SELECT * FROM users
WHERE name = $1 limit 1;

-- name: GetUser :one
select * from users
where id = $1 limit 1;

-- name: ResetUsers :exec
delete from users;

-- name: GetUsers :many
select * from users;