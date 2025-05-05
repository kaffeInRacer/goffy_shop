CREATE TABLE IF NOT EXISTS payment (
    id             BIGSERIAL PRIMARY KEY,
    order_id       BIGINT NOT NULL
        REFERENCES orders(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    payment_method INTEGER NOT NULL
        REFERENCES payment_method(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    amount         DOUBLE PRECISION NOT NULL,
    status         INTEGER    NOT NULL,
    created_time    TIMESTAMP  NOT NULL DEFAULT now(),
    updated_time    TIMESTAMP  NOT NULL DEFAULT now()
);