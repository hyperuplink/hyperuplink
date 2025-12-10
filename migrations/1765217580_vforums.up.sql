CREATE MATERIALIZED VIEW vforums AS
SELECT
    f.*,
    (SELECT COUNT(*)
     FROM topics t
     WHERE t.forum_id = f.id) AS topics,

    (SELECT COUNT(*)
     FROM replies r
     JOIN topics t ON r.topic_id = t.id
     WHERE t.forum_id = f.id) AS replies,

    (SELECT MAX(r.created_at)
     FROM replies r
     JOIN topics t ON r.topic_id = t.id
     WHERE t.forum_id = f.id) AS "last_reply_at"
FROM forums f;
