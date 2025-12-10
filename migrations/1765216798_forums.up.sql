CREATE TABLE forums (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(32),
  slug VARCHAR(32),
  position SMALLINT DEFAULT 0,
  category_id UUID NOT NULL,
  description VARCHAR(128),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_forums_slug ON forums (
  slug, category_id, deleted_at
);

CREATE INDEX idx_forums_categories_id ON forums (category_id);
