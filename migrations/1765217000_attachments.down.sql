DROP INDEX IF EXISTS idx_attachments_users_id;
DROP INDEX IF EXISTS idx_attachments_checksum_author;

DROP TABLE IF EXISTS attachments;

DROP FUNCTION IF EXISTS check_attachment_ids() CASCADE;
