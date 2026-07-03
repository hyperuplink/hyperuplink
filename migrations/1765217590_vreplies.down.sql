DROP INDEX IF EXISTS idx_vreplies_fts;

DROP TRIGGER IF EXISTS refresh_vrepliess ON replies;
DROP TRIGGER IF EXISTS refresh_vrepliess ON users;
DROP FUNCTION IF EXISTS refresh_vrepliess();
DROP MATERIALIZED VIEW IF EXISTS vrepliess;
