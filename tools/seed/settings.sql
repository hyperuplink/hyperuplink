-- Applied before the users are created, because the address type decides
-- whether `email_is_jid` means anything at all: with the default "email only"
-- a JID is validated as an email address, quietly stored as one, and the user
-- you meant to reach over XMPP is then an email user.
--
-- Nothing here sets the theme or the base URL: the first is a matter of taste
-- and the second depends on where the board is being run, so both belong to
-- whoever is doing the seeding.

BEGIN;

UPDATE settings
SET json_value = json_value || '{
  "allowed_address_type": 2
}'::jsonb
WHERE id = 'auth';

UPDATE settings
SET json_value = json_value || '{
  "enable_about": true,
  "about": "This board runs Hyperuplink.",
  "enable_contact": true,
  "contact": "Reach the sysop at sysop@example.org.",
  "enable_terms": true,
  "terms": "Be excellent to each other.",
  "enable_privacy_policy": true,
  "privacy_policy": "This board stores what you post, and nothing besides."
}'::jsonb
WHERE id = 'general';

UPDATE settings
SET json_value = json_value || '{
  "allow_kind_poll": true
}'::jsonb
WHERE id = 'topics';

UPDATE settings
SET json_value = json_value || '{
  "enable_attachments": true,
  "storage_provider_id": "local-storage",
  "max_size": 4194304
}'::jsonb
WHERE id = 'attachments';

UPDATE settings
SET json_value = json_value || '{
  "enable_picture": true,
  "picture_storage_provider_id": "local-storage"
}'::jsonb
WHERE id = 'profiles';

COMMIT;
