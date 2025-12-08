CREATE TABLE permissions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  role VARCHAR(16) DEFAULT 'user',
  unit VARCHAR(32) DEFAULT NULL,
  forum_id UUID NOT NULL,
  bits BIT(3) DEFAULT B'000',
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT fk_unit FOREIGN KEY (unit) REFERENCES units (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_forum FOREIGN KEY (forum_id) REFERENCES forums (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_permissions_units_id ON permissions (unit);
CREATE INDEX idx_permissions_forums_id ON permissions (forum_id);
