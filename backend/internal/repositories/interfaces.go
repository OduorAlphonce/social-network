package repositories

import (
	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

// UserRepository stores and reads user account records.
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUserProfile(id uuid.UUID) (*models.User, error)
	DeleteUser(id uuid.UUID) error
}

// SessionRepository stores and reads authenticated session records.
type SessionRepository interface {
	CreateSession(session *models.Session) error
	GetSessionByID(id uuid.UUID) (*models.Session, error)
	DeleteSession(id uuid.UUID) error
}

// FollowersRepository stores and reads user follow relationships.
type FollowersRepository interface {
	Follow(followerID, followeeID uuid.UUID, status models.Status) error
	Unfollow(followerID, followeeID uuid.UUID) error
	AcceptFollower(followerID, followeeID uuid.UUID) error
	RejectFollower(followerID, followeeID uuid.UUID) error
	GetFollowers(userID uuid.UUID) ([]*models.User, error)
	GetFollowing(userID uuid.UUID) ([]*models.User, error)
	GetStatus(followerID, followeeID uuid.UUID) (models.Status, error)
}

// PostRepository stores posts and returns post rows hydrated with viewer state.
type PostRepository interface {
	CreatePost(post *models.Post) error
	GetPostByID(id, viewerID uuid.UUID) (*models.PostWithAuthor, error)
	ListPosts(query models.PostQuery, viewerID uuid.UUID) ([]*models.PostWithAuthor, error)
	ListHomeFeed(viewerID uuid.UUID, limit, offset int) ([]*models.PostWithAuthor, error)
	ListProfilePosts(profileUserID, viewerID uuid.UUID, limit, offset int) ([]*models.PostWithAuthor, error)
	ListGroupFeed(groupID, viewerID uuid.UUID, limit, offset int) ([]*models.PostWithAuthor, error)
}

// CommentRepository stores comments and returns comment rows hydrated with viewer state.
type CommentRepository interface {
	CreateComment(comment *models.Comment) error
	GetCommentByID(id, viewerID uuid.UUID) (*models.CommentWithAuthor, error)
	ListCommentTreeByPost(postID, viewerID uuid.UUID, limit, offset int) ([]*models.CommentWithAuthor, error)
}

// PostAudienceRepository stores selected-follower audiences for private posts.
type PostAudienceRepository interface {
	ReplacePostAudience(postID uuid.UUID, userIDs []uuid.UUID) error
	ListPostAudience(postID uuid.UUID) ([]uuid.UUID, error)
	IsPostAudienceMember(postID, userID uuid.UUID) (bool, error)
}

// PostVoteRepository stores mutually exclusive post votes and vote summaries.
type PostVoteRepository interface {
	SetPostVote(postID, userID uuid.UUID, vote models.VoteValue) (*models.VoteSummary, error)
	DeletePostVote(postID, userID uuid.UUID) (*models.VoteSummary, error)
	GetPostVoteSummary(postID, viewerID uuid.UUID) (*models.VoteSummary, error)
}

// CommentVoteRepository stores mutually exclusive comment votes and vote summaries.
type CommentVoteRepository interface {
	SetCommentVote(commentID, userID uuid.UUID, vote models.VoteValue) (*models.VoteSummary, error)
	DeleteCommentVote(commentID, userID uuid.UUID) (*models.VoteSummary, error)
	GetCommentVoteSummary(commentID, viewerID uuid.UUID) (*models.VoteSummary, error)
}

// GroupMembershipRepository reads group membership state needed by scoped features.
type GroupMembershipRepository interface {
	IsAcceptedGroupMember(groupID, userID uuid.UUID) (bool, error)
}
