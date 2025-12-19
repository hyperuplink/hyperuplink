CREATE TABLE topics (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  short_id TEXT NOT NULL,
  name VARCHAR(78),
  slug VARCHAR(78),
  forum_id UUID NOT NULL,
  author_id UUID NOT NULL,
  kind VARCHAR(16) DEFAULT 'regular',
  anonymous BOOLEAN DEFAULT FALSE,
  pinned BOOLEAN DEFAULT FALSE,
  text TEXT NOT NULL,
  html TEXT NOT NULL,
  poll_options VARCHAR(78) [],
  views BIGINT DEFAULT 0,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  ended_at TIMESTAMP DEFAULT NULL,
  moderated_at TIMESTAMP DEFAULT NULL,
  spammed_at TIMESTAMP DEFAULT NULL,
  locked_at TIMESTAMP DEFAULT NULL,
  deleted_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT fk_forum FOREIGN KEY (forum_id) REFERENCES forums (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES users (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_topics_slug ON topics (
  slug, forum_id
) WHERE deleted_at IS NULL;

CREATE INDEX idx_topics_short_id ON topics (short_id);
CREATE INDEX idx_topics_forums_id ON topics (forum_id);
CREATE INDEX idx_topics_users_id ON topics (author_id);
