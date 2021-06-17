
-- +goose Up
create extension if not exists "uuid-ossp";

create table if not exists bobpos.product_category (
     category_id uuid unique not null,
     category_name varchar(200) unique not null,
     created_at timestamp with time zone not null default current_timestamp,
     updated_at timestamp with time zone not null default current_timestamp,
     primary key (category_name)
);

insert into
    bobpos.product_category(category_id, category_name)
    VALUES (uuid_generate_v4(), 'CHICKEN'),
           (uuid_generate_v4(), 'COW'),
           (uuid_generate_v4(), 'TURKEY'),
           (uuid_generate_v4(), 'OTHER') on conflict do nothing;
alter table bobpos.products add foreign key (category) references bobpos.product_category;

-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
drop table if exists bobpos.product_category cascade;
drop extension if exists "uuid-ossp";
alter table bobpos.products drop constraint if exists products_category_fkey;

-- SQL section 'Down' is executed when this migration is rolled back

