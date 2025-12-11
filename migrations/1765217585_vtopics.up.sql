CREATE MATERIALIZED VIEW vtopics AS
SELECT
    t.*,
    COALESCE((SELECT COUNT(*)
      FROM replies r
      WHERE r.topic_id = t.id), 0) AS replies,

    (SELECT MAX(r.created_at)
     FROM replies r
     WHERE r.topic_id = t.id) AS "last_reply_at",

    u.username AS author,
    u.email AS author_email,
    u.profile_picture AS author_profile_picture,
    u.signature AS author_signature,

    c.Name AS category_name,
    c.Slug AS category_slug,

    f.Name AS forum_name,
    f.Slug AS forum_slug
FROM topics t
LEFT JOIN users u ON u.id = t.author_id
LEFT JOIN forums f ON f.id = t.forum_id
LEFT JOIN categories c ON c.id = f.category_id;
