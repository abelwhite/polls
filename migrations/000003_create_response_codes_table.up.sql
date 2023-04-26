-- Filename: migrations/000003_create_response_codes_table.up.sql
CREATE TABLE IF NOT EXISTS response_codes (
  response_code_id bigserial PRIMARY KEY,
  code text NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);