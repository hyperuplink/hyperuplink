CREATE TABLE attachments (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  author_id UUID NOT NULL,
  filename TEXT NOT NULL,
  mime_type TEXT NOT NULL,
  checksum TEXT NOT NULL,
  on_upload_hook_output TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT now(),
  moderated_at TIMESTAMP DEFAULT NULL,
  spammed_at TIMESTAMP DEFAULT NULL,
  deleted_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES users (
    id
  ) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_attachments_checksum_author ON attachments (
  checksum, author_id
) WHERE deleted_at IS NULL;

CREATE INDEX idx_attachments_users_id ON attachments (author_id);

CREATE OR REPLACE FUNCTION check_attachment_ids() RETURNS TRIGGER AS $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM unnest(NEW.attachment_ids) AS attachment_id
    WHERE NOT EXISTS (
      SELECT 1 FROM attachments WHERE id = attachment_id
    )
  ) THEN
    RAISE EXCEPTION 'attachment_ids references a non-existent attachment'
      USING ERRCODE = 'foreign_key_violation';
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
