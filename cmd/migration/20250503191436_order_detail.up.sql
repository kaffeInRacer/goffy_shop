CREATE TABLE IF NOT EXISTS order_detail (
 id BIGSERIAL PRIMARY KEY,
 product_id ulid NOT NULL
    REFERENCES product(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,
 quantity INTEGER NOT NULL
);
