CREATE TABLE postevents (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  type VARCHAR(16) NOT NULL, -- "flag", "pollvote", "answervote", "rsvpresponse"
  author_id UUID NOT NULL,
  target VARCHAR(16) NOT NULL,
  topic_id UUID DEFAULT NULL,
  reply_id UUID DEFAULT NULL,
  selection SMALLINT DEFAULT -1,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES users (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_postevents_type ON postevents (
  type, author_id, target, topic_id, reply_id
) WHERE deleted_at IS NULL;

CREATE INDEX idx_postevents_type ON postevents (type);
CREATE INDEX idx_postevents_author_id ON postevents (author_id);
CREATE INDEX idx_postevents_target ON postevents (target);
CREATE INDEX idx_postevents_topic_id ON postevents (topic_id);
CREATE INDEX idx_postevents_reply_id ON postevents (reply_id);
