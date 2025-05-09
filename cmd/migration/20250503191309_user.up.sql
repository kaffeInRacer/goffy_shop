CREATE EXTENSION IF NOT EXISTS ulid;
CREATE TYPE gender_type AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS users (
  id           ulid NOT NULL DEFAULT gen_ulid() PRIMARY KEY,
  username     VARCHAR(255)    NOT NULL UNIQUE,
  slug         varchar(255)    NOT NULL UNIQUE,
  email        VARCHAR(255)    NOT NULL UNIQUE,
  password     VARCHAR(255)    NOT NULL,
  role         VARCHAR(50)     NOT NULL,
  gender       gender_type     NOT NULL,
  created_at   TIMESTAMP       DEFAULT NOW(),
  updated_at   TIMESTAMP       DEFAULT NOW()
);
