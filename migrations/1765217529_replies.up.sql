CREATE TABLE replies (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  short_id TEXT NOT NULL,
  topic_id UUID NOT NULL,
  reply_id UUID DEFAULT NULL,
  author_id UUID NOT NULL,
  text TEXT NOT NULL,
  html TEXT NOT NULL,
  attachment_ids UUID [] DEFAULT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  moderated_at TIMESTAMP DEFAULT NULL,
  spammed_at TIMESTAMP DEFAULT NULL,
  deleted_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT fk_topic FOREIGN KEY (topic_id) REFERENCES topics (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_reply FOREIGN KEY (reply_id) REFERENCES replies (
    id
  ) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES users (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_replies_short_id ON replies (short_id);
CREATE INDEX idx_replies_topics_id ON replies (topic_id);
CREATE INDEX idx_replies_replies_id ON replies (reply_id);
CREATE INDEX idx_replies_users_id ON replies (author_id);

CREATE TRIGGER trg_replies_attachment_ids
  BEFORE INSERT OR UPDATE OF attachment_ids ON replies
  FOR EACH ROW EXECUTE FUNCTION check_attachment_ids();
