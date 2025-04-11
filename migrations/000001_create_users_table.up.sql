CREATE TYPE roles AS ENUM (
  'admin',
  'staff',
  'viewer'
);

CREATE TABLE IF NOT EXISTS users (
  user_id BigSerial PRIMARY KEY,
  username varchar(100) UNIQUE NOT NULL,
  userpassword text NOT NULL,
  role roles NOT NULL,
  created_at timestamp DEFAULT now()
);
