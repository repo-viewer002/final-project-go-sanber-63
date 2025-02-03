-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE genres (
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
CREATE TRIGGER genres_modified_at_trigger BEFORE
UPDATE ON genres FOR EACH ROW EXECUTE FUNCTION update_modified_at();
-- +migrate StatementEnd

-- +migrate Up
-- +migrate StatementBegin
INSERT INTO genres (name, description, created_by, modified_by)
VALUES (
    'mistery',
    'this is a mistery genre',
    'system',
    'system'
  ),
  (
    'drama',
    'this is a drama genre',
    'system',
    'system'
  ),
  (
    'romance',
    'this is a romance genre',
    'system',
    'system'
  ),
  (
    'sci-fi',
    'this is a sci-fi genre',
    'system',
    'system'
  ),
  (
    'fantasy',
    'this is a fiction genre',
    'system',
    'system'
  ),
  (
    'thriller',
    'this is a thriller genre',
    'system',
    'system'
  ),
  (
    'education',
    'this is a education genre',
    'system',
    'system'
  ),
  (
    'action/adventure',
    'this is a action/adventure genre',
    'system',
    'system'
  ),
  (
    'personal-growth',
    'this is a personal-growth genre',
    'system',
    'system'
  ),
  (
    'biography',
    'this is a biography genre',
    'system',
    'system'
  ),
  (
    'historical',
    'this is a historical genre',
    'system',
    'system'
  ),
  (
    'travel',
    'this is a travel genre',
    'system',
    'system'
  );
-- +migrate StatementEnd