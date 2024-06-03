-- +goose Up
create table users (
    id          bigserial primary key,
    username    varchar(255) not null,
    password    varchar(255) not null,
    role        int not null,
    updated_at  timestamptz default now(),
    created_at  timestamptz default now(),
    deleted_at  timestamptz default null
);

insert into users(username, password, role) values ('Seller','$2a$12$yT.dJTZnu4FRJq9zXw0mBOA/xmZHJPVi5ni13Zk9Pn6E0QmwKkZTu',1);
insert into users(username, password, role) values ('Buyer','$2a$12$yT.dJTZnu4FRJq9zXw0mBOA/xmZHJPVi5ni13Zk9Pn6E0QmwKkZTu',2);

-- +goose Down
drop table users;