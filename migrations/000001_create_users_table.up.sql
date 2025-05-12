CREATE TYPE roles AS ENUM ('admin', 'staff');

CREATE TABLE IF NOT EXISTS users (
  user_id BIGSERIAL PRIMARY KEY,
  username VARCHAR(100) UNIQUE NOT NULL,
  email CITEXT UNIQUE NOT NULL,
  password_hash BYTEA NOT NULL,
  role roles NOT NULL DEFAULT 'admin',
  activated BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMP DEFAULT now()
);