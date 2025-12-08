CREATE TABLE forums (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(32),
  slug VARCHAR(32),
  category_id UUID NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_forums_categories_id ON forums (category_id);
