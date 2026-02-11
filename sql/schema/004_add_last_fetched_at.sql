-- +goose Up
alter table feeds
add column if not exists last_fetched_at timestamp;

-- +goose Down
alter table feeds
drop column last_fecthed_at;