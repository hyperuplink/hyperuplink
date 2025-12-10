CREATE MATERIALIZED VIEW vtopics AS
SELECT
    t.*,
    COALESCE((SELECT COUNT(*)
      FROM replies r
      WHERE r.topic_id = t.id), 0) AS replies,

    (SELECT MAX(r.created_at)
     FROM replies r
     WHERE r.topic_id = t.id) AS "last_reply_at"
FROM topics t;
