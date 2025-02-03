-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE borrowed_books (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  borrow_id UUID NOT NULL,
  book_id UUID NOT NULL,
  FOREIGN KEY (borrow_id) REFERENCES borrows(id) ON DELETE CASCADE,
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);
-- +migrate StatementEnd