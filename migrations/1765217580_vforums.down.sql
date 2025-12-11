DROP TRIGGER IF EXISTS refresh_vforums ON replies;
DROP TRIGGER IF EXISTS refresh_vforums ON topics;
DROP TRIGGER IF EXISTS refresh_vforums ON forums;
DROP TRIGGER IF EXISTS refresh_vforums ON categories;
DROP FUNCTION IF EXISTS refresh_vforums();
DROP MATERIALIZED VIEW IF EXISTS vforums;
