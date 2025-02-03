-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE book_genres (
  book_id UUID,
  genre_id UUID,
  PRIMARY KEY (book_id, genre_id),
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE
);
-- +migrate StatementEnd