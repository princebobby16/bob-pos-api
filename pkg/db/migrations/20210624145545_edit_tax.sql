
-- +goose Up
drop table if exists bobpos.tax;
create table if not exists bobpos.tax (
    id uuid unique not null,
    name varchar(200) not null,
    tax_rate double precision not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    primary key (tax_rate)
);
alter table bobpos.products add foreign key (tax) references bobpos.tax;
-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
alter table bobpos.products drop constraint products_tax_fkey;
drop table if exists bobpos.tax;
-- SQL section 'Down' is executed when this migration is rolled back

