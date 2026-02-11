-- +goose Up

create table feeds(
    id uuid primary key,
    name varchar unique not null,
    url varchar unique not null,
    user_id uuid not null ,
    foreign key (user_id)
    references users(id) on delete cascade
);

-- +goose Down
drop table feeds;