CREATE TABLE IF NOT EXISTS product (
  id           ulid NOT NULL DEFAULT gen_ulid() PRIMARY KEY ,
  name         VARCHAR(255)   NOT NULL,
  description  TEXT,
  price        DOUBLE PRECISION NOT NULL,
  stock        INTEGER        NOT NULL DEFAULT 0,
  category_id  INTEGER        NOT NULL
      REFERENCES product_category(id)
      ON UPDATE CASCADE
      ON DELETE RESTRICT,
  created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);