-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) UNIQUE NOT NULL,
  description VARCHAR(255) NOT NULL,
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
CREATE TRIGGER roles_modified_at_trigger BEFORE
UPDATE ON roles FOR EACH ROW EXECUTE FUNCTION update_modified_at();
-- +migrate StatementEnd
-- +migrate Up

-- +migrate StatementBegin
INSERT INTO roles (name, description, created_by, modified_by)
VALUES (
    'admin',
    'this is admin description',
    'system',
    'system'
  ),
  (
    'librarian',
    'this is librarian description',
    'system',
    'system'
  ),
  (
    'member',
    'this is member description',
    'system',
    'system'
  );
-- +migrate StatementEnd