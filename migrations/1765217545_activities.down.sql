DROP INDEX IF EXISTS idx_activities_created_at;
DROP INDEX IF EXISTS idx_activities_subject;
DROP INDEX IF EXISTS idx_activities_kind;
DROP INDEX IF EXISTS idx_activities_actor_id;
DROP INDEX IF EXISTS idx_activities_dedupe;

DROP TABLE IF EXISTS activities;
