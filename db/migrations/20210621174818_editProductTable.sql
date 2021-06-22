
-- +goose Up
alter table bobpos.products add column barcode varchar(300) unique not null;
-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
alter table bobpos.products drop column if exists barcode;
-- SQL section 'Down' is executed when this migration is rolled back

