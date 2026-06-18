package services

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

func TestPermissionServiceIsFollowing(t *testing.T) {
	followers := newFakeFollowersRepository()
	users := newFakeUserRepository()
	groups := newFakeGroupMembershipRepository()
	permissionSvc := NewPermissionService(followers, users, groups)

	followerID := uuid.Must(uuid.NewV4())
	followeeID := uuid.Must(uuid.NewV4())

	// No follow relationship yet
	isFollowing, err := permissionSvc.IsFollowing(context.Background(), followerID, followeeID)
	if err != nil {
		t.Fatalf("IsFollowing error: %v", err)
	}
	if isFollowing {
		t.Error("expected isFollowing to be false when no status exists")
	}

	// Pending relationship
	followers.status[followerKey{followerID: followerID, followeeID: followeeID}] = models.Pending
	isFollowing, err = permissionSvc.IsFollowing(context.Background(), followerID, followeeID)
	if err != nil {
		t.Fatalf("IsFollowing error: %v", err)
	}
	if isFollowing {
		t.Error("expected isFollowing to be false for pending status")
	}

	// Accepted relationship
	followers.status[followerKey{followerID: followerID, followeeID: followeeID}] = models.Accepted
	isFollowing, err = permissionSvc.IsFollowing(context.Background(), followerID, followeeID)
	if err != nil {
		t.Fatalf("IsFollowing error: %v", err)
	}
	if !isFollowing {
		t.Error("expected isFollowing to be true for accepted status")
	}
}

func TestPermissionServiceCanViewProfile(t *testing.T) {
	followers := newFakeFollowersRepository()
	users := newFakeUserRepository()
	groups := newFakeGroupMembershipRepository()
	permissionSvc := NewPermissionService(followers, users, groups)

	viewerID := uuid.Must(uuid.NewV4())
	profileID := uuid.Must(uuid.NewV4())

	// Case 1: Viewing own profile
	canView, err := permissionSvc.CanViewProfile(context.Background(), viewerID, viewerID)
	if err != nil {
		t.Fatalf("CanViewProfile error: %v", err)
	}
	if !canView {
		t.Error("expected owner to be able to view own profile")
	}

	// Case 2: Public profile
	publicUser := &models.User{
		ID:       profileID,
		Email:    "public@example.com",
		IsPublic: true,
	}
	users.add(publicUser)

	canView, err = permissionSvc.CanViewProfile(context.Background(), viewerID, profileID)
	if err != nil {
		t.Fatalf("CanViewProfile error: %v", err)
	}
	if !canView {
		t.Error("expected viewer to be able to view public profile")
	}

	// Case 3: Private profile (no relationship)
	privateID := uuid.Must(uuid.NewV4())
	privateUser := &models.User{
		ID:       privateID,
		Email:    "private@example.com",
		IsPublic: false,
	}
	users.add(privateUser)

	canView, err = permissionSvc.CanViewProfile(context.Background(), viewerID, privateID)
	if err != nil {
		t.Fatalf("CanViewProfile error: %v", err)
	}
	if canView {
		t.Error("expected viewer to be blocked from viewing private profile with no follow status")
	}

	// Case 4: Private profile (pending follow)
	followers.status[followerKey{followerID: viewerID, followeeID: privateID}] = models.Pending
	canView, err = permissionSvc.CanViewProfile(context.Background(), viewerID, privateID)
	if err != nil {
		t.Fatalf("CanViewProfile error: %v", err)
	}
	if canView {
		t.Error("expected viewer to be blocked from viewing private profile when follow is pending")
	}

	// Case 5: Private profile (accepted follow)
	followers.status[followerKey{followerID: viewerID, followeeID: privateID}] = models.Accepted
	canView, err = permissionSvc.CanViewProfile(context.Background(), viewerID, privateID)
	if err != nil {
		t.Fatalf("CanViewProfile error: %v", err)
	}
	if !canView {
		t.Error("expected accepted follower to be able to view private profile")
	}
}

func TestPermissionServiceIsGroupMember(t *testing.T) {
	followers := newFakeFollowersRepository()
	users := newFakeUserRepository()
	groups := newFakeGroupMembershipRepository()
	permissionSvc := NewPermissionService(followers, users, groups)

	groupID := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())

	// Non-member
	isMember, err := permissionSvc.IsGroupMember(context.Background(), groupID, userID)
	if err != nil {
		t.Fatalf("IsGroupMember error: %v", err)
	}
	if isMember {
		t.Error("expected isMember to be false for non-member")
	}

	// Accepted member
	groups.accepted[groupMemberKey{groupID: groupID, userID: userID}] = true
	isMember, err = permissionSvc.IsGroupMember(context.Background(), groupID, userID)
	if err != nil {
		t.Fatalf("IsGroupMember error: %v", err)
	}
	if !isMember {
		t.Error("expected isMember to be true for accepted member")
	}
}

func TestPermissionServiceCanViewPost(t *testing.T) {
	followers := newFakeFollowersRepository()
	users := newFakeUserRepository()
	groups := newFakeGroupMembershipRepository()
	permissionSvc := NewPermissionService(followers, users, groups)

	authorID := uuid.Must(uuid.NewV4())
	viewerID := uuid.Must(uuid.NewV4())
	viewerStr := viewerID.String()

	// Case 1: Public Post
	publicPost := &models.Post{
		ID:      uuid.Must(uuid.NewV4()),
		UserID:  &authorID,
		Privacy: models.PostPrivacyPublic,
	}
	canView, err := permissionSvc.CanViewPost(context.Background(), &viewerStr, publicPost)
	if err != nil {
		t.Fatalf("CanViewPost error: %v", err)
	}
	if !canView {
		t.Error("expected public post to be viewable by anyone")
	}

	// Case 2: Almost Private Post (Not following)
	almostPrivatePost := &models.Post{
		ID:      uuid.Must(uuid.NewV4()),
		UserID:  &authorID,
		Privacy: models.PostPrivacyAlmostPrivate,
	}
	canView, err = permissionSvc.CanViewPost(context.Background(), &viewerStr, almostPrivatePost)
	if err != nil {
		t.Fatalf("CanViewPost error: %v", err)
	}
	if canView {
		t.Error("expected almost_private post to be blocked for non-followers")
	}

	// Case 3: Almost Private Post (Accepted follow)
	followers.status[followerKey{followerID: viewerID, followeeID: authorID}] = models.Accepted
	canView, err = permissionSvc.CanViewPost(context.Background(), &viewerStr, almostPrivatePost)
	if err != nil {
		t.Fatalf("CanViewPost error: %v", err)
	}
	if !canView {
		t.Error("expected almost_private post to be viewable by accepted followers")
	}

	// Case 4: Owner views private post
	privatePost := &models.Post{
		ID:      uuid.Must(uuid.NewV4()),
		UserID:  &authorID,
		Privacy: models.PostPrivacyPrivate,
	}
	authorStr := authorID.String()
	canView, err = permissionSvc.CanViewPost(context.Background(), &authorStr, privatePost)
	if err != nil {
		t.Fatalf("CanViewPost error: %v", err)
	}
	if !canView {
		t.Error("expected post owner to be able to view their own private post")
	}

	// Case 5: Non-owner views private post
	canView, err = permissionSvc.CanViewPost(context.Background(), &viewerStr, privatePost)
	if err != nil {
		t.Fatalf("CanViewPost error: %v", err)
	}
	if canView {
		t.Error("expected private post to be blocked for non-owners")
	}

	// Case 6: Group Post (Non-member)
	groupID := uuid.Must(uuid.NewV4())
	groupPost := &models.Post{
		ID:      uuid.Must(uuid.NewV4()),
		UserID:  &authorID,
		GroupID: &groupID,
		Privacy: models.PostPrivacyPublic,
	}
	canView, err = permissionSvc.CanViewPost(context.Background(), &viewerStr, groupPost)
	if err != nil {
		t.Fatalf("CanViewPost error: %v", err)
	}
	if canView {
		t.Error("expected group post to be blocked for non-group-members")
	}

	// Case 7: Group Post (Accepted member)
	groups.accepted[groupMemberKey{groupID: groupID, userID: viewerID}] = true
	canView, err = permissionSvc.CanViewPost(context.Background(), &viewerStr, groupPost)
	if err != nil {
		t.Fatalf("CanViewPost error: %v", err)
	}
	if !canView {
		t.Error("expected group post to be viewable by group members")
	}
}
