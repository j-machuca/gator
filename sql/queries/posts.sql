-- name: InsertPost :one
insert into posts(
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
) Values (
    $1, $2, $3, $4, $5, $6, $7, $8
) returning *;


-- name: GetPosts :many
select title,description,published_at,url from posts limit $1;