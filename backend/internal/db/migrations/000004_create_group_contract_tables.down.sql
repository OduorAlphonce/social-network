DROP INDEX IF EXISTS idx_group_members_group_status;
DROP INDEX IF EXISTS idx_group_members_user_status;
DROP TABLE IF EXISTS group_members;

DROP INDEX IF EXISTS idx_groups_created_at;
DROP INDEX IF EXISTS idx_groups_creator_created;
DROP TABLE IF EXISTS groups;
