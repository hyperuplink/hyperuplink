DROP TRIGGER IF EXISTS refresh_vtopics ON replies;
DROP TRIGGER IF EXISTS refresh_vtopics ON topics;
DROP TRIGGER IF EXISTS refresh_vtopics ON forums;
DROP TRIGGER IF EXISTS refresh_vtopics ON categories;
DROP TRIGGER IF EXISTS refresh_vtopics ON users;
DROP FUNCTION IF EXISTS refresh_vtopics();
DROP MATERIALIZED VIEW IF EXISTS vtopics;
