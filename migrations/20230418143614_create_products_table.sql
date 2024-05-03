-- +goose Up
CREATE TABLE if not exists products(
    id INT PRIMARY KEY NOT NULL ,
    brand TEXT NOT NULL,
    category TEXT NOT NULL,
    quantity INT NOT NULL,
    price FLOAT NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

-- +goose Down
DROP TABLE if exists products;