-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE borrows (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  borrowed_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  return_deadline TIMESTAMP,
  returned_time TIMESTAMP,
  status VARCHAR(20) DEFAULT 'borrowed' CHECK (status IN ('borrowed', 'returned', 'overdue')),
  created_by VARCHAR(255) NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +migrate StatementEnd