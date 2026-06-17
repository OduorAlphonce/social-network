DROP INDEX IF EXISTS idx_comment_votes_vote;
DROP INDEX IF EXISTS idx_comment_votes_user;
DROP TABLE IF EXISTS comment_votes;

DROP INDEX IF EXISTS idx_post_votes_vote;
DROP INDEX IF EXISTS idx_post_votes_user;
DROP TABLE IF EXISTS post_votes;

DROP INDEX IF EXISTS idx_comments_user_created;
DROP INDEX IF EXISTS idx_comments_parent_created;
DROP INDEX IF EXISTS idx_comments_post_created;
DROP TABLE IF EXISTS comments;

DROP INDEX IF EXISTS idx_post_audiences_user;
DROP TABLE IF EXISTS post_audiences;

DROP INDEX IF EXISTS idx_posts_privacy_created;
DROP INDEX IF EXISTS idx_posts_group_created;
DROP INDEX IF EXISTS idx_posts_user_created;
DROP INDEX IF EXISTS idx_posts_created_at;
DROP TABLE IF EXISTS posts;
