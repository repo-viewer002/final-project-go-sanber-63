-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE books (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  description VARCHAR(255) NOT NULL,
  authors VARCHAR(255),
  publisher VARCHAR(255),
  publish_year INTEGER,
  stock INTEGER DEFAULT 1,
  borrowed INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  created_by VARCHAR(255) NOT NULL,
  modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  modified_by VARCHAR(255) NOT NULL
);
-- +migrate StatementEnd

-- +migrate Up
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_modified_at() RETURNS TRIGGER AS $$ BEGIN NEW.modified_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER books_modified_at_trigger BEFORE
UPDATE ON books FOR EACH ROW EXECUTE FUNCTION update_modified_at();
-- +migrate StatementEnd