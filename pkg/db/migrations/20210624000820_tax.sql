
-- +goose Up
create table if not exists bobpos.tax (
    id uuid unique not null,
    name varchar(200) not null,
    tax_rate varchar(200) not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    primary key (id)
);
-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
drop table if exists bobpos.tax;
-- SQL section 'Down' is executed when this migration is rolled back

