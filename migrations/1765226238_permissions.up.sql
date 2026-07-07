CREATE TABLE permissions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  group_id VARCHAR(32) DEFAULT NULL,
  category_id UUID DEFAULT NULL,
  bits BIT(3) DEFAULT B'000',
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES groups (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX idx_permissions_group_category ON permissions (
  COALESCE(group_id, ''),
  COALESCE(category_id, '00000000-0000-0000-0000-000000000000'::UUID)
) WHERE deleted_at IS NULL;

CREATE INDEX idx_permissions_group_id ON permissions (group_id);
CREATE INDEX idx_permissions_category_id ON permissions (category_id);
