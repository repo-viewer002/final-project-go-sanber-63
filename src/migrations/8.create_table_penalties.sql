-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE penalties (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  borrow_id UUID NOT NULL,
  total_amount INTEGER NOT NULL CHECK (total_amount >= 0),
  paid_amount INTEGER,
  paid_off_time TIMESTAMP,
  status VARCHAR(20) DEFAULT 'unpaid' CHECK (status IN ('unpaid', 'installment', 'paid')),
  FOREIGN KEY (borrow_id) REFERENCES borrows(id) ON DELETE CASCADE
);
-- +migrate StatementEnd