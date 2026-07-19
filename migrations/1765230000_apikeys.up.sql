CREATE TABLE apikeys (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  name VARCHAR(64) DEFAULT '',
  secret_hash VARCHAR(64) NOT NULL,
  last_used_at TIMESTAMP DEFAULT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_apikeys_secret_hash ON apikeys (
  secret_hash
) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_apikeys_user_id ON apikeys (
  user_id
);
