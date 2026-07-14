CREATE TABLE activities (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  kind VARCHAR(32) NOT NULL,
  actor_id UUID NOT NULL,
  subject VARCHAR(16) NOT NULL,
  subject_id UUID DEFAULT NULL,
  dedupe_key TEXT DEFAULT NULL,
  count INTEGER DEFAULT 1,
  context JSONB DEFAULT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT fk_actor FOREIGN KEY (actor_id) REFERENCES users (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_activities_dedupe ON activities (
  kind, actor_id, dedupe_key
) WHERE dedupe_key IS NOT NULL AND deleted_at IS NULL;

CREATE INDEX idx_activities_actor_id ON activities (actor_id);
CREATE INDEX idx_activities_kind ON activities (kind);
CREATE INDEX idx_activities_subject ON activities (subject, subject_id);
CREATE INDEX idx_activities_created_at ON activities (created_at);
