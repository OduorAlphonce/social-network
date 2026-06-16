package repositories

import (
	"database/sql"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	_ "github.com/mattn/go-sqlite3"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

func TestPostRepositoryReturnsViewerVoteInReadModel(t *testing.T) {
	db := newPostCommentTestDB(t)
	ids := seedPostCommentTestRows(t, db)

	repo := NewPostRepository(db)
	post, err := repo.GetPostByID(ids.postID, ids.viewerID)
	if err != nil {
		t.Fatalf("GetPostByID returned error: %v", err)
	}

	if post.Post.ID != ids.postID {
		t.Fatalf("post id = %s, want %s", post.Post.ID, ids.postID)
	}
	if post.Author == nil || post.Author.ID != ids.authorID {
		t.Fatalf("author = %#v, want %s", post.Author, ids.authorID)
	}
	if post.ViewerVote != models.ViewerVoteLike {
		t.Fatalf("viewer vote = %q, want %q", post.ViewerVote, models.ViewerVoteLike)
	}
}

func TestCommentRepositoryReturnsTreeWithViewerVotes(t *testing.T) {
	db := newPostCommentTestDB(t)
	ids := seedPostCommentTestRows(t, db)

	repo := NewCommentRepository(db)
	comments, err := repo.ListCommentTreeByPost(ids.postID, ids.viewerID, 10, 0)
	if err != nil {
		t.Fatalf("ListCommentTreeByPost returned error: %v", err)
	}
	if len(comments) != 2 {
		t.Fatalf("got %d comments, want 2", len(comments))
	}

	byID := map[uuid.UUID]*models.CommentWithAuthor{}
	for _, comment := range comments {
		byID[comment.Comment.ID] = comment
	}
	if byID[ids.parentCommentID].ViewerVote != models.ViewerVoteNone {
		t.Fatalf("parent viewer vote = %q, want none", byID[ids.parentCommentID].ViewerVote)
	}
	if byID[ids.replyCommentID].ViewerVote != models.ViewerVoteDislike {
		t.Fatalf("reply viewer vote = %q, want dislike", byID[ids.replyCommentID].ViewerVote)
	}
	if byID[ids.replyCommentID].Comment.ParentCommentID == nil || *byID[ids.replyCommentID].Comment.ParentCommentID != ids.parentCommentID {
		t.Fatalf("reply parent = %#v, want %s", byID[ids.replyCommentID].Comment.ParentCommentID, ids.parentCommentID)
	}
}

func TestPostAudienceRepositoryReplacesAndChecksMembership(t *testing.T) {
	db := newPostCommentTestDB(t)
	ids := seedPostCommentTestRows(t, db)
	memberID := ids.otherViewerID

	repo := NewPostAudienceRepository(db)
	if err := repo.ReplacePostAudience(ids.postID, []uuid.UUID{ids.viewerID}); err != nil {
		t.Fatalf("ReplacePostAudience returned error: %v", err)
	}
	if err := repo.ReplacePostAudience(ids.postID, []uuid.UUID{memberID}); err != nil {
		t.Fatalf("ReplacePostAudience replace returned error: %v", err)
	}

	members, err := repo.ListPostAudience(ids.postID)
	if err != nil {
		t.Fatalf("ListPostAudience returned error: %v", err)
	}
	if len(members) != 1 || members[0] != memberID {
		t.Fatalf("members = %#v, want only %s", members, memberID)
	}

	isMember, err := repo.IsPostAudienceMember(ids.postID, memberID)
	if err != nil {
		t.Fatalf("IsPostAudienceMember returned error: %v", err)
	}
	if !isMember {
		t.Fatal("expected replacement member to be in audience")
	}
	oldMember, err := repo.IsPostAudienceMember(ids.postID, ids.viewerID)
	if err != nil {
		t.Fatalf("IsPostAudienceMember old member returned error: %v", err)
	}
	if oldMember {
		t.Fatal("expected previous member to be removed")
	}
}

func TestVoteRepositoriesUpdateCountsAndViewerVote(t *testing.T) {
	db := newPostCommentTestDB(t)
	ids := seedPostCommentTestRows(t, db)
	postVotes := NewPostVoteRepository(db)
	commentVotes := NewCommentVoteRepository(db)

	postSummary, err := postVotes.SetPostVote(ids.postID, ids.otherViewerID, models.VoteValueDislike)
	if err != nil {
		t.Fatalf("SetPostVote returned error: %v", err)
	}
	if postSummary.LikeCount != 1 || postSummary.DislikeCount != 1 || postSummary.ViewerVote != models.ViewerVoteDislike {
		t.Fatalf("post summary = %#v, want 1 like, 1 dislike, viewer dislike", postSummary)
	}

	postSummary, err = postVotes.DeletePostVote(ids.postID, ids.viewerID)
	if err != nil {
		t.Fatalf("DeletePostVote returned error: %v", err)
	}
	if postSummary.LikeCount != 0 || postSummary.DislikeCount != 1 || postSummary.ViewerVote != models.ViewerVoteNone {
		t.Fatalf("post summary after delete = %#v, want 0 likes, 1 dislike, viewer none", postSummary)
	}

	commentSummary, err := commentVotes.SetCommentVote(ids.replyCommentID, ids.viewerID, models.VoteValueLike)
	if err != nil {
		t.Fatalf("SetCommentVote returned error: %v", err)
	}
	if commentSummary.LikeCount != 1 || commentSummary.DislikeCount != 0 || commentSummary.ViewerVote != models.ViewerVoteLike {
		t.Fatalf("comment summary = %#v, want 1 like, 0 dislikes, viewer like", commentSummary)
	}
}

type postCommentSeedIDs struct {
	authorID        uuid.UUID
	viewerID        uuid.UUID
	otherViewerID   uuid.UUID
	postID          uuid.UUID
	parentCommentID uuid.UUID
	replyCommentID  uuid.UUID
}

func newPostCommentTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}

	schema := `
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			dob TEXT NOT NULL,
			avatar TEXT,
			nickname TEXT,
			about_me TEXT,
			is_public INTEGER NOT NULL DEFAULT 1,
			follower_count INTEGER NOT NULL DEFAULT 0,
			following_count INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL
		);
		CREATE TABLE posts (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			group_id TEXT,
			content TEXT,
			image_url TEXT,
			privacy TEXT NOT NULL CHECK (privacy IN ('public', 'almost_private', 'private')),
			comment_count INTEGER NOT NULL DEFAULT 0,
			like_count INTEGER NOT NULL DEFAULT 0,
			dislike_count INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT,
			deleted_at TEXT,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		);
		CREATE TABLE post_audiences (
			post_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			PRIMARY KEY (post_id, user_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE post_votes (
			post_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			vote TEXT NOT NULL CHECK (vote IN ('like', 'dislike')),
			created_at TEXT NOT NULL,
			updated_at TEXT,
			PRIMARY KEY (post_id, user_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE comments (
			id TEXT PRIMARY KEY,
			post_id TEXT NOT NULL,
			user_id TEXT,
			parent_comment_id TEXT,
			content TEXT,
			image_url TEXT,
			like_count INTEGER NOT NULL DEFAULT 0,
			dislike_count INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			deleted_at TEXT,
			updated_at TEXT,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
			FOREIGN KEY (parent_comment_id) REFERENCES comments(id)
		);
		CREATE TABLE comment_votes (
			comment_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			vote TEXT NOT NULL CHECK (vote IN ('like', 'dislike')),
			created_at TEXT NOT NULL,
			updated_at TEXT,
			PRIMARY KEY (comment_id, user_id),
			FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	return db
}

func seedPostCommentTestRows(t *testing.T, db *sql.DB) postCommentSeedIDs {
	t.Helper()

	ids := postCommentSeedIDs{
		authorID:        uuid.Must(uuid.FromString("6f5d9a18-5c4f-4b7a-9e9a-7a5d2efc44b1")),
		viewerID:        uuid.Must(uuid.FromString("0dd6e443-0998-4f50-a4cf-1a40a0536213")),
		otherViewerID:   uuid.Must(uuid.FromString("d54271d8-e248-4260-acf2-b6f1264830af")),
		postID:          uuid.Must(uuid.FromString("8a6de4c1-2f4a-4c52-a38f-61d76c9e7d11")),
		parentCommentID: uuid.Must(uuid.FromString("623b1e60-babc-44db-a4a1-0c3d071a80a3")),
		replyCommentID:  uuid.Must(uuid.FromString("a1ac5c54-6cb6-41bc-a88f-55ad4fb3a6d2")),
	}
	now := time.Date(2026, 6, 16, 8, 15, 0, 0, time.UTC)

	insertUser(t, db, ids.authorID, "amina@example.com", "Amina", "Njeri")
	insertUser(t, db, ids.viewerID, "viewer@example.com", "Vera", "Viewer")
	insertUser(t, db, ids.otherViewerID, "other@example.com", "Otis", "Other")

	_, err := db.Exec(`
		INSERT INTO posts (
			id, user_id, content, image_url, privacy,
			comment_count, like_count, dislike_count, created_at
		)
		VALUES (?, ?, ?, ?, 'public', 1, 1, 0, ?)
	`, ids.postID.String(), ids.authorID.String(), "First hike", "/uploads/posts/hike.gif", now.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("insert post: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO post_votes (post_id, user_id, vote, created_at)
		VALUES (?, ?, 'like', ?)
	`, ids.postID.String(), ids.viewerID.String(), now.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("insert post vote: %v", err)
	}

	deletedAt := now.Add(5 * time.Minute)
	_, err = db.Exec(`
		INSERT INTO comments (
			id, post_id, user_id, content, like_count, dislike_count, created_at, deleted_at
		)
		VALUES (?, ?, ?, ?, 0, 0, ?, ?)
	`, ids.parentCommentID.String(), ids.postID.String(), ids.authorID.String(), "deleted retained text", now.Add(time.Minute).Format(time.RFC3339), deletedAt.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("insert parent comment: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO comments (
			id, post_id, user_id, parent_comment_id, content, like_count, dislike_count, created_at
		)
		VALUES (?, ?, ?, ?, ?, 0, 1, ?)
	`, ids.replyCommentID.String(), ids.postID.String(), ids.authorID.String(), ids.parentCommentID.String(), "That view is unreal.", now.Add(2*time.Minute).Format(time.RFC3339))
	if err != nil {
		t.Fatalf("insert reply comment: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO comment_votes (comment_id, user_id, vote, created_at)
		VALUES (?, ?, 'dislike', ?)
	`, ids.replyCommentID.String(), ids.viewerID.String(), now.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("insert comment vote: %v", err)
	}

	return ids
}

func insertUser(t *testing.T, db *sql.DB, id uuid.UUID, email, firstName, lastName string) {
	t.Helper()
	_, err := db.Exec(`
		INSERT INTO users (
			id, email, password_hash, first_name, last_name, dob,
			avatar, nickname, about_me, is_public, created_at
		)
		VALUES (?, ?, 'hash', ?, ?, '1998-04-12', '/uploads/avatars/user.png', 'nick', NULL, 1, ?)
	`, id.String(), email, firstName, lastName, time.Date(2026, 6, 16, 8, 0, 0, 0, time.UTC).Format(time.RFC3339))
	if err != nil {
		t.Fatalf("insert user %s: %v", id, err)
	}
}
