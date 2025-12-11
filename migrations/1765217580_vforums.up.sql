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

    (SELECT MAX(r.created_at)
     FROM replies r
     JOIN topics t ON r.topic_id = t.id
     WHERE t.forum_id = f.id) AS "last_reply_at",

    c.name AS category_name,
    c.slug AS category_slug
FROM forums f
LEFT JOIN categories c ON c.id = f.category_id;
