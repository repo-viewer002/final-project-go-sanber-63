-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE penalty_payments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  penalty_id UUID NOT NULL,
  amount INTEGER NOT NULL,
  paid_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (penalty_id) REFERENCES penalties(id) ON DELETE CASCADE
);
-- +migrate StatementEnd