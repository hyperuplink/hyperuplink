CREATE MATERIALIZED VIEW vreplies AS
SELECT
    r.*,
    u.username AS author_username,
    u.role AS author_role,
    u.member_of AS author_member_of,
    u.email AS author_email,
    u.profile_picture AS author_profile_picture,
    u.signature AS author_signature,
    u.created_at AS author_joined_at
FROM replies r
LEFT JOIN users u ON u.id = r.author_id;

CREATE OR REPLACE FUNCTION refresh_vreplies() RETURNS TRIGGER
  AS $$
  BEGIN
      REFRESH MATERIALIZED VIEW vreplies;
    RETURN NULL;
  END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER refresh_vreplies
 AFTER INSERT OR UPDATE ON users
 FOR EACH STATEMENT EXECUTE PROCEDURE refresh_vreplies();

CREATE TRIGGER refresh_vreplies
 AFTER INSERT ON replies
 FOR EACH STATEMENT EXECUTE PROCEDURE refresh_vreplies();

