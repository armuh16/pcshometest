-- +goose Up
create table transactions (
    id           bigserial primary key,
    buyer_id     int not null,
    seller_id    int not null,
    origin       varchar(255) not null,
    destination  varchar(255) not null,
    items        json not null,
    grand_total  float not null,
    status       int not null default 1,
    coupons      int not null default 0,
    updated_at   timestamptz default now(),
    created_at   timestamptz default now(),
    deleted_at   timestamptz default null,
    foreign key  (buyer_id) references users (id),
    foreign key  (seller_id) references users (id)
);

-- +goose Down
drop table transactions;