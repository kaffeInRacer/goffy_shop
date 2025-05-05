CREATE TABLE IF NOT EXISTS payment_method (
  id     SERIAL PRIMARY KEY,
  name   VARCHAR(100) NOT NULL,
  status INTEGER      NOT NULL
);