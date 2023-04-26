CREATE TABLE IF NOT EXISTS users (
  users_id bigserial PRIMARY KEY,
  username text NOT NULL,
  email citext UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  activated boolean NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);