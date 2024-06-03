-- +goose Up
create table products (
    id          bigserial primary key,
    name        varchar(255) not null,
    description varchar(255) not null,
    price       float not null,
    seller_id   int not null,
    updated_at timestamptz default now(),
    created_at timestamptz default now(),
    deleted_at timestamptz default null,
    foreign key (seller_id) references users (id)
);

-- +goose Down
drop table products;