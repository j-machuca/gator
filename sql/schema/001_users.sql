-- +goose Up 

create table users(
    id uuid primary key,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name varchar UNIQUE not null
);


-- +goose Down
drop table users;