-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id                     UUID PRIMARY KEY,
    coupon_code            VARCHAR(255),
    created_at             BIGINT NOT NULL
);

CREATE TABLE products (
    id                     UUID PRIMARY KEY,
    name                   VARCHAR(255) NOT NULL,
    category               VARCHAR(255),
    price                  FLOAT NOT NULL,
    created_at             BIGINT NOT NULL
);

CREATE TABLE order_product (
  id                    UUID PRIMARY KEY,
  order_id              UUID NOT NULL,
  product_id            UUID NOT NULL,
  FOREIGN KEY (order_id) REFERENCES orders(id),
  FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE coupons (
    id                     VARCHAR(255) PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_product;
DROP TABLE orders;
DROP TABLE products;
DROP TABLE coupons;
-- +goose StatementEnd
