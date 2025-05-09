CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  username varchar(20) NOT NULL,
  age int,
  first_name varchar(50) NOT NULL,
  last_name varchar(50) NOT NULL,
  email citext UNIQUE NOT NULL,
  hashed_password bytea NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
