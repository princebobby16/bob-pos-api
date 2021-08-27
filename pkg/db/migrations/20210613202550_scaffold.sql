
-- +goose Up
create schema bobpos;
create table if not exists bobpos.products(
    id uuid not null unique,
    name varchar not null,
    category varchar,
    weight varchar,
    cost_price float not null,
    tax float not null,
    profit_margin float not null,
    image bytea,
    number_in_stock integer not null,
    created_at       timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    primary key (id)
);
-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
drop table if exists bobpos.products;
drop schema bobpos;
-- SQL section 'Down' is executed when this migration is rolled back

