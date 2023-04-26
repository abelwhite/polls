-- Filename: migrations/000005_create_question_responses_table.up.sql
CREATE TABLE question_responses (
  question_response_id bigserial PRIMARY KEY,
  question_id BIGINT REFERENCES questions(question_id),
  response_code_id BIGINT REFERENCES response_codes(response_code_id),
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

