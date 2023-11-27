create table if not exists orders (
    order_uid varchar primary key,
    order_json JSONB
);