CREATE MATERIALIZED VIEW vtopics AS
SELECT
    t.*,
    COALESCE((SELECT COUNT(*)
      FROM replies r
      WHERE r.topic_id = t.id), 0) AS replies,

    (SELECT MAX(r.created_at)
     FROM replies r
     WHERE r.topic_id = t.id) AS "last_reply_at",

    COALESCE((SELECT COUNT(*)
      FROM activities a
      WHERE a.kind = 'topic_view' AND a.subject = 'topic'
        AND a.subject_id = t.id AND a.deleted_at IS NULL), 0) AS views,

    u.username AS author_username,
    u.role AS author_role,
    u.member_of AS author_member_of,
    u.email AS author_email,
    u.profile_picture AS author_profile_picture,
    u.signature_text AS author_signature_text,
    u.signature_html AS author_signature_html,
    u.created_at AS author_joined_at,

    c.Name AS category_name,
    c.Slug AS category_slug,

    f.Name AS forum_name,
    f.Slug AS forum_slug
FROM topics t
LEFT JOIN users u ON u.id = t.author_id
LEFT JOIN forums f ON f.id = t.forum_id
LEFT JOIN categories c ON c.id = f.category_id;

CREATE OR REPLACE FUNCTION refresh_vtopics() RETURNS TRIGGER
  AS $$
  BEGIN
      REFRESH MATERIALIZED VIEW vtopics;
    RETURN NULL;
  END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER refresh_vtopics
 AFTER INSERT OR UPDATE ON users
 FOR EACH STATEMENT EXECUTE PROCEDURE refresh_vtopics();

CREATE TRIGGER refresh_vtopics
 AFTER INSERT OR UPDATE ON categories
 FOR EACH STATEMENT EXECUTE PROCEDURE refresh_vtopics();

CREATE TRIGGER refresh_vtopics
 AFTER INSERT OR UPDATE ON forums
 FOR EACH STATEMENT EXECUTE PROCEDURE refresh_vtopics();

CREATE TRIGGER refresh_vtopics
 AFTER INSERT ON topics
 FOR EACH STATEMENT EXECUTE PROCEDURE refresh_vtopics();

CREATE TRIGGER refresh_vtopics
 AFTER INSERT ON replies
 FOR EACH STATEMENT EXECUTE PROCEDURE refresh_vtopics();

-- We could enable this to refresh vtopics every time activities changes, in
-- case a "topic_view" event was added. However, this could lead to **a lot** of
-- refreshing. Updates to the views field aren't usually *that* important, so
-- it is probably okay to not trigger a refresh for every new activity.
--
-- CREATE TRIGGER refresh_vtopics
--  AFTER INSERT ON activities
--  FOR EACH STATEMENT EXECUTE PROCEDURE refresh_vtopics();

-- IMPORTANT: the tsvector expression below MUST stay identical to the one used
-- by the search query in services/repositories/topic/topic.vsearch.go,
-- otherwise the planner cannot use this index.
CREATE INDEX IF NOT EXISTS idx_vtopics_fts ON vtopics
USING gin (
  to_tsvector('english', coalesce(name, '') || ' ' || coalesce(text, ''))
);

