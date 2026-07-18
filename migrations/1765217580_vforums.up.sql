CREATE MATERIALIZED VIEW vforums AS
SELECT
  f.*,
  COALESCE((SELECT COUNT(*)
    FROM topics t
    WHERE t.forum_id = f.id), 0) AS topics,

  COALESCE((SELECT COUNT(*)
    FROM replies r
    JOIN topics t ON r.topic_id = t.id
    WHERE t.forum_id = f.id), 0) AS replies,

  GREATEST(
    (SELECT MAX(r.created_at)
     FROM replies r
     JOIN topics t ON r.topic_id = t.id
     WHERE t.forum_id = f.id),
    (SELECT MAX(t.created_at)
     FROM topics t
     WHERE t.forum_id = f.id)
  ) AS "last_reply_at",

  c.name AS category_name,
  c.slug AS category_slug
FROM forums f
LEFT JOIN categories c ON c.id = f.category_id;

CREATE OR REPLACE FUNCTION refresh_vforums() RETURNS TRIGGER
  AS $$
  BEGIN
      REFRESH MATERIALIZED VIEW vforums;
    RETURN NULL;
  END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER refresh_vforums
 AFTER INSERT OR UPDATE ON categories
 FOR EACH STATEMENT EXECUTE PROCEDURE refresh_vforums();

CREATE TRIGGER refresh_vforums
 AFTER INSERT OR UPDATE ON forums
 FOR EACH STATEMENT EXECUTE PROCEDURE refresh_vforums();

CREATE TRIGGER refresh_vforums
 AFTER INSERT OR UPDATE ON topics
 FOR EACH STATEMENT EXECUTE PROCEDURE refresh_vforums();

CREATE TRIGGER refresh_vforums
 AFTER INSERT OR UPDATE ON replies
 FOR EACH STATEMENT EXECUTE PROCEDURE refresh_vforums();
