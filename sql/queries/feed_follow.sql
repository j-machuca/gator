-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    INSERT INTO feed_follows(id,created_at,updated_at,user_id,feed_id) Values ($1,$2,$3,$4,$5)
    returning *
) SELECT 
    iff.*,
    f.name as feed,
    u.name as username
    from inserted_feed_follow as iff
    inner join feeds as f
    on iff.feed_id = f.id
    inner join users as u
    on iff.user_id = u.id;


-- name: GetFeedFollowsForUser :many
select u.name as username, f.name as feed
from feed_follows as ff
inner join feeds as f
on ff.feed_id = f.id
inner join users as u
on ff.user_id = u.id
where ff.user_id = $1;


-- name: ResetFeedFollowings :exec
delete from feed_follows;

-- name: Unfollow :exec
DELETE FROM feed_follows ff
USING feeds f, users u
WHERE ff.feed_id = f.id
  AND ff.user_id = u.id
  AND f.url = $1
  AND u.name = $2;