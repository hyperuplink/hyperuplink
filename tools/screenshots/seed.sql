-- Applied on top of tools/seed/, which is shared with the development board.
-- Only the two things the screenshots need and a development board does not
-- belong here: the theme the manual shows throughout, and the port this
-- particular board is reachable on.

BEGIN;

UPDATE settings
SET json_value = json_value
  || '{"theme":"macos9","colorscheme":"hyperuplink-light"}'::jsonb
WHERE id = 'theme';

UPDATE settings
SET json_value = json_value
  || '{"name":"Hyperuplink","base_url":"http://localhost:3101"}'::jsonb
WHERE id = 'system';

COMMIT;
