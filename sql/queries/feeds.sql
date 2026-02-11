-- name: CreateFeed :one
INSERT INTO FEEDS(id,name,url,user_id) values(
    $1,
    $2,
    $3,
    $4
)
returning *;

-- name: ResetFeed :exec

delete from feeds;


-- name: GetFeeds :many

SELECT f.name as name,f.url as url,u.name as username from
feeds as f
left join users as u
on f.user_id = u.id;


-- name: GetFeedByUrl :one

select * from feeds where url=$1;


-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at = $2,
updated_at = $2
where feeds.id = $1;

-- name: GetNextFeedToFetch :one

select id,url,min(last_fecthed_at) as last_fecthed_at from feeds
group by id
order by last_fecthed_at asc nulls first;