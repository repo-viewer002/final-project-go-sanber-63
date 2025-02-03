-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(100) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  email VARCHAR(100) UNIQUE,
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  address VARCHAR(255),
  phone_number VARCHAR(50),
  is_penalized BOOLEAN DEFAULT FALSE,
  penalty_duration TIMESTAMP,
  role_id UUID,
  status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'deactivated', 'suspended')),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  created_by VARCHAR(255) NOT NULL,
  modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  modified_by VARCHAR(255) NOT NULL,
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE
  SET NULL
);
-- +migrate StatementEnd

-- +migrate Up
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_modified_at() RETURNS TRIGGER AS $$ BEGIN NEW.modified_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER users_modified_at_trigger BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_modified_at();
-- +migrate StatementEnd

-- Testing Purpose

-- +migrate Up
-- +migrate StatementBegin
INSERT INTO users (
  username, 
  password, 
  email, 
  first_name, 
  last_name, 
  address, 
  phone_number, 
  is_penalized, 
  role_id, 
  status, 
  created_by, 
  modified_by
) 
VALUES (
  'main_admin', 
  '$2a$14$hXwNaqrbWEkwJLF/8nKCoONhzL87ayUFb9DiKRA1CCo1eCesfhYeO', -- Hashed password: "Admin@123"
  'main_admin@mail.com', 
  'Main', 
  'Admin', 
  'Unknown Location', 
  '08123456789', 
  FALSE, 
  (SELECT id FROM roles WHERE name = 'admin'),
  'active', 
  'system',
  'system'
)
ON CONFLICT (username) DO NOTHING;
-- +migrate StatementEnd
