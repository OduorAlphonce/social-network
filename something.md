Based on the context provided, I need to extract the follow/user component code that was mentioned as already provided in previous messages. However, I don't see those actual component files in the current conversation - I only see the existing files like Friends.jsx, Profile.jsx, etc.

Looking at what you have, it seems you want to continue building the follow system and user cards. Since you mentioned the files were provided in previous messages, let me create the missing components that would complete the follow system integration based on the existing code patterns in your codebase.

Here are the missing components needed to complete the follow system:

1. src/components/user/UserCard.jsx
```jsx
import { useNavigate } from "react-router";
import avatarFallback from "../../assets/user.svg";
import "../../styles/user-card.css";

const toDisplayName = (user) =>
  `${user.first_name || ""} ${user.last_name || ""}`.trim() ||
  user.nickname ||
  "Unnamed User";

const UserCard = ({ user, actions, onClick }) => {
  const navigate = useNavigate();

  const handleClick = () => {
    if (onClick) {
      onClick(user);
    } else {
      navigate(`/user/${user.id}`);
    }
  };

  return (
    <div className="user-card" onClick={handleClick}>
      <img
        src={user.avatar || avatarFallback}
        alt={`${toDisplayName(user)}'s avatar`}
        className="user-card__avatar"
      />
      <div className="user-card__info">
        <div className="user-card__name">{toDisplayName(user)}</div>
        {user.nickname && (
          <div className="user-card__nickname">@{user.nickname}</div>
        )}
        {user.about_me && (
          <div className="user-card__bio">{user.about_me}</div>
        )}
        <div className="user-card__badge">
          {user.is_public ? (
            <span className="user-card__badge-public">Public</span>
          ) : (
            <span className="user-card__badge-private">Private</span>
          )}
        </div>
      </div>
      {actions && <div className="user-card__actions">{actions}</div>}
    </div>
  );
};

export default UserCard;
```
2. src/styles/user-card.css
```css
.user-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background-color: var(--bg-elevated);
  border: 1px solid var(--card-border);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.user-card:hover {
  border-color: var(--hyperpop-violet);
  box-shadow: var(--glow-violet);
}

.user-card__avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
  background-color: var(--bg-surface);
}

.user-card__info {
  flex: 1;
  min-width: 0;
}

.user-card__name {
  color: var(--text-primary);
  font-weight: 600;
  font-size: 0.95rem;
}

.user-card__nickname {
  color: var(--text-muted);
  font-size: 0.85rem;
}

.user-card__bio {
  color: var(--text-secondary);
  font-size: 0.85rem;
  margin-top: 4px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.user-card__badge {
  margin-top: 4px;
}

.user-card__badge-public,
.user-card__badge-private {
  font-size: 0.7rem;
  padding: 2px 8px;
  border-radius: 4px;
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.user-card__badge-public {
  background-color: rgba(204, 255, 0, 0.15);
  color: var(--acid-green);
}

.user-card__badge-private {
  background-color: rgba(255, 51, 102, 0.15);
  color: var(--punch-pink);
}

.user-card__actions {
  flex-shrink: 0;
}
```
3. src/components/follow/FollowAction.jsx
```jsx
import { useState } from "react";
import { apiFetch } from "../../utils/api";
import "../../styles/follow-action.css";

const FollowAction = ({ targetUserId, initialStatus = "unfollowed", isPrivate = false, onStatusChange }) => {
  const [status, setStatus] = useState(initialStatus); // "unfollowed" | "following" | "requested"
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleClick = async () => {
    setLoading(true);
    setError("");

    try {
      if (status === "unfollowed") {
        // Send follow request
        const body = { following_id: targetUserId };
        await apiFetch("/api/followers/follow", {
          method: "POST",
          body,
        });
        const newStatus = isPrivate ? "requested" : "following";
        setStatus(newStatus);
        onStatusChange?.(newStatus);
      } else if (status === "following") {
        // Unfollow
        await apiFetch("/api/followers/unfollow", {
          method: "POST",
          body: { following_id: targetUserId },
        });
        setStatus("unfollowed");
        onStatusChange?.("unfollowed");
      } else if (status === "requested") {
        // Cancel follow request
        await apiFetch("/api/followers/cancel", {
          method: "POST",
          body: { following_id: targetUserId },
        });
        setStatus("unfollowed");
        onStatusChange?.("unfollowed");
      }
    } catch (err) {
      setError(err.message || "Failed to update follow status");
    } finally {
      setLoading(false);
    }
  };

  const getButtonConfig = () => {
    switch (status) {
      case "following":
        return {
          label: loading ? "..." : "Following",
          className: "follow-action--following",
          hoverLabel: "Unfollow",
        };
      case "requested":
        return {
          label: loading ? "..." : "Requested",
          className: "follow-action--requested",
          hoverLabel: "Cancel Request",
        };
      default:
        return {
          label: loading ? "..." : "Follow",
          className: "follow-action--follow",
          hoverLabel: null,
        };
    }
  };

  const config = getButtonConfig();

  return (
    <div className="follow-action">
      <button
        className={`follow-action__btn ${config.className}`}
        onClick={handleClick}
        disabled={loading}
        data-hover-label={config.hoverLabel}
      >
        {loading ? "..." : config.label}
      </button>
      {error && <div className="follow-action__error">{error}</div>}
    </div>
  );
};

export default FollowAction;
```
4. src/styles/follow-action.css
```css
.follow-action__btn {
  padding: 8px 20px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
  min-width: 100px;
}

.follow-action--follow {
  background: var(--hyperpop-violet);
  color: white;
}

.follow-action--follow:hover {
  box-shadow: 0 0 20px rgba(139, 92, 246, 0.4);
  transform: translateY(-1px);
}

.follow-action--following {
  background: transparent;
  color: var(--text-primary);
  border: 1px solid var(--text-muted);
}

.follow-action--following:hover {
  border-color: var(--punch-pink);
  color: var(--punch-pink);
}

.follow-action--following:hover::after {
  content: attr(data-hover-label);
}

.follow-action--following:hover .follow-action__btn-text {
  display: none;
}

.follow-action--requested {
  background: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--text-muted);
}

.follow-action--requested:hover {
  border-color: var(--punch-pink);
  color: var(--punch-pink);
}

.follow-action--requested:hover::after {
  content: attr(data-hover-label);
}

.follow-action__btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.follow-action__error {
  color: var(--punch-pink);
  font-size: 0.8rem;
  margin-top: 4px;
}
```
5. src/components/follow/FollowRequestCard.jsx
```jsx
import { useState } from "react";
import { apiFetch } from "../../utils/api";
import avatarFallback from "../../assets/user.svg";
import "../../styles/follow-request-card.css";

const toDisplayName = (user) =>
  `${user.first_name || ""} ${user.last_name || ""}`.trim() ||
  user.nickname ||
  "Unnamed User";

const FollowRequestCard = ({ request, onAccept, onReject }) => {
  const [loading, setLoading] = useState(false);

  const handleAccept = async () => {
    setLoading(true);
    try {
      await apiFetch("/api/followers/accept", {
        method: "POST",
        body: { follower_id: request.id },
      });
      onAccept?.(request.id);
    } catch (err) {
      console.error("Failed to accept request:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleReject = async () => {
    setLoading(true);
    try {
      await apiFetch("/api/followers/reject", {
        method: "POST",
        body: { follower_id: request.id },
      });
      onReject?.(request.id);
    } catch (err) {
      console.error("Failed to reject request:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="follow-request-card">
      <div className="follow-request-card__user">
        <img
          src={request.avatar || avatarFallback}
          alt={`${toDisplayName(request)}'s avatar`}
          className="follow-request-card__avatar"
        />
        <div>
          <div className="follow-request-card__name">
            {toDisplayName(request)}
          </div>
          {request.nickname && (
            <div className="follow-request-card__nickname">
              @{request.nickname}
            </div>
          )}
        </div>
      </div>
      <div className="follow-request-card__actions">
        <button
          className="follow-request-card__accept"
          onClick={handleAccept}
          disabled={loading}
        >
          Accept
        </button>
        <button
          className="follow-request-card__reject"
          onClick={handleReject}
          disabled={loading}
        >
          Reject
        </button>
      </div>
    </div>
  );
};

export default FollowRequestCard;
```
6. src/styles/follow-request-card.css
```css
.follow-request-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background-color: var(--bg-elevated);
  border: 1px solid var(--card-border);
  border-radius: 12px;
  transition: all 0.2s ease;
}

.follow-request-card:hover {
  border-color: var(--hyperpop-violet);
}

.follow-request-card__user {
  display: flex;
  align-items: center;
  gap: 12px;
}

.follow-request-card__avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.follow-request-card__name {
  color: var(--text-primary);
  font-weight: 500;
}

.follow-request-card__nickname {
  color: var(--text-muted);
  font-size: 0.85rem;
}

.follow-request-card__actions {
  display: flex;
  gap: 8px;
}

.follow-request-card__accept,
.follow-request-card__reject {
  padding: 6px 16px;
  border-radius: 6px;
  font-weight: 600;
  font-size: 0.85rem;
  cursor: pointer;
  border: none;
  transition: all 0.2s ease;
}

.follow-request-card__accept {
  background: var(--hyperpop-violet);
  color: white;
}

.follow-request-card__accept:hover {
  box-shadow: 0 0 20px rgba(139, 92, 246, 0.3);
}

.follow-request-card__reject {
  background: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--text-muted);
}

.follow-request-card__reject:hover {
  border-color: var(--punch-pink);
  color: var(--punch-pink);
}

.follow-request-card__accept:disabled,
.follow-request-card__reject:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
```
7. src/components/follow/FollowRequestsList.jsx
```jsx
import { useEffect, useState } from "react";
import { apiFetch } from "../../utils/api";
import FollowRequestCard from "./FollowRequestCard";
import "../../styles/follow-requests-list.css";

const FollowRequestsList = ({ onRequestCountChange }) => {
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const fetchRequests = async () => {
    setLoading(true);
    setError("");
    try {
      const data = await apiFetch("/api/followers/pending");
      setRequests(Array.isArray(data) ? data : []);
      onRequestCountChange?.(data.length);
    } catch (err) {
      setError(err.message || "Failed to load follow requests");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRequests();
  }, []);

  const handleAccept = (requestId) => {
    setRequests((prev) => prev.filter((r) => r.id !== requestId));
    onRequestCountChange?.(requests.length - 1);
  };

  const handleReject = (requestId) => {
    setRequests((prev) => prev.filter((r) => r.id !== requestId));
    onRequestCountChange?.(requests.length - 1);
  };

  if (loading) {
    return (
      <div className="follow-requests-list">
        <div className="follow-requests-list__skeleton" />
        <div className="follow-requests-list__skeleton" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="follow-requests-list__error">{error}</div>
    );
  }

  if (requests.length === 0) {
    return (
      <div className="follow-requests-list__empty">
        No pending follow requests.
      </div>
    );
  }

  return (
    <div className="follow-requests-list">
      {requests.map((request) => (
        <FollowRequestCard
          key={request.id}
          request={request}
          onAccept={handleAccept}
          onReject={handleReject}
        />
      ))}
    </div>
  );
};

export default FollowRequestsList;
```
8. src/styles/follow-requests-list.css
```css
.follow-requests-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.follow-requests-list__skeleton {
  height: 72px;
  background: linear-gradient(
    90deg,
    var(--bg-surface) 25%,
    var(--bg-elevated) 50%,
    var(--bg-surface) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 12px;
}

@keyframes shimmer {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}

.follow-requests-list__error {
  color: var(--punch-pink);
  text-align: center;
  padding: 20px;
}

.follow-requests-list__empty {
  color: var(--text-muted);
  text-align: center;
  padding: 20px;
}
```

9. Updated src/pages/Friends.jsx
```jsx
import { useEffect, useState } from "react";
import { useAuth } from "../context/useAuth";
import { apiFetch, ApiError } from "../utils/api";
import FollowRequestsList from "../components/follow/FollowRequestsList";
import UserCard from "../components/user/UserCard";
import FollowAction from "../components/follow/FollowAction";
import avatarFallback from "../assets/user.svg";
import "../styles/friends.css";

const toDisplayName = (user) =>
  `${user.first_name || ""} ${user.last_name || ""}`.trim() ||
  user.nickname ||
  "Unnamed User";

const Friends = () => {
  const { currentUser } = useAuth();
  const [query, setQuery] = useState("");
  const [suggestedUsers, setSuggestedUsers] = useState([]);
  const [pendingRequests, setPendingRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState({});
  const [actionMessage, setActionMessage] = useState("");
  const [error, setError] = useState("");
  const [requestCount, setRequestCount] = useState(0);

  useEffect(() => {
    if (!currentUser) return;

    const loadPage = async () => {
      setLoading(true);
      setError("");

      try {
        const [pending, suggestions] = await Promise.all([
          apiFetch("/api/followers/pending"),
          apiFetch("/api/users/search"),
        ]);
        setPendingRequests(Array.isArray(pending) ? pending : pending || []);
        setRequestCount(pending?.length || 0);
        setSuggestedUsers(
          Array.isArray(suggestions) ? suggestions : suggestions || []
        );
      } catch (err) {
        setError(
          err instanceof ApiError ? err.message : "Unable to load friends data."
        );
      } finally {
        setLoading(false);
      }
    };

    loadPage();
  }, [currentUser]);

  const refreshPending = async () => {
    try {
      const pending = await apiFetch("/api/followers/pending");
      setPendingRequests(Array.isArray(pending) ? pending : pending || []);
      setRequestCount(pending?.length || 0);
    } catch {
      // don't block the user if pending refresh fails
    }
  };

  const loadSuggestions = async (searchTerm = "") => {
    setError("");
    try {
      const results = await apiFetch(
        `/api/users/search${searchTerm ? `?query=${encodeURIComponent(searchTerm)}` : ""}`
      );
      setSuggestedUsers(Array.isArray(results) ? results : results || []);
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Unable to load suggestions."
      );
    }
  };

  const handleSearchSubmit = async (event) => {
    event.preventDefault();
    await loadSuggestions(query);
  };

  const handleFollowStatusChange = (userId, newStatus) => {
    // Update the suggested users list to reflect new status
    setSuggestedUsers((prev) =>
      prev.map((user) => {
        if (user.id === userId) {
          // We don't store follow status in the user object directly,
          // but we could add a `_followStatus` field if needed
          return { ...user, _followStatus: newStatus };
        }
        return user;
      })
    );
    setActionMessage(
      newStatus === "following"
        ? "Followed successfully!"
        : newStatus === "unfollowed"
        ? "Unfollowed successfully!"
        : "Follow request sent!"
    );
  };

  if (!currentUser) {
    return <div className="friends-page">Loading friends...</div>;
  }

  return (
    <div className="friends-page">
      <div className="friends-header">
        <h1>Friends</h1>
        <p>
          Find public profiles, send friend requests, and accept pending follow
          requests.
        </p>
      </div>

      {error && (
        <div className="profile-state profile-state--error">{error}</div>
      )}
      {actionMessage && <div className="profile-state">{actionMessage}</div>}

      <section className="friends-section">
        <div className="friends-section__header">
          <h2>
            Pending requests{" "}
            {requestCount > 0 && (
              <span className="friends-badge">{requestCount}</span>
            )}
          </h2>
        </div>

        <FollowRequestsList onRequestCountChange={setRequestCount} />
      </section>

      <section className="friends-section">
        <div className="friends-section__header friends-section__header--search">
          <h2>Discover public profiles</h2>
          <form className="friends-search" onSubmit={handleSearchSubmit}>
            <input
              type="search"
              value={query}
              placeholder="Search by username, first name, or last name"
              onChange={(event) => setQuery(event.target.value)}
              className="friends-search__input"
            />
            <button type="submit" className="profile-btn profile-btn--primary">
              Search
            </button>
          </form>
        </div>

        {loading ? (
          <div className="profile-skeleton profile-skeleton--row" />
        ) : suggestedUsers.length === 0 ? (
          <div className="profile-state">
            No public profiles found. Try a different search term.
          </div>
        ) : (
          <div className="friends-list">
            {suggestedUsers.map((user) => (
              <UserCard
                key={user.id}
                user={user}
                actions={
                  <FollowAction
                    targetUserId={user.id}
                    initialStatus="unfollowed"
                    isPrivate={!user.is_public}
                    onStatusChange={(status) =>
                      handleFollowStatusChange(user.id, status)
                    }
                  />
                }
              />
            ))}
          </div>
        )}
      </section>
    </div>
  );
};

export default Friends;
```
10. Updated src/pages/UserProfile.jsx
```jsx
import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router";
import { useAuth } from "../context/useAuth";
import { apiFetch, ApiError } from "../utils/api";
import FollowAction from "../components/follow/FollowAction";
import Post from "../components/Post";
import avatarFallback from "../assets/user.svg";
import "../styles/profile.css";

const toDisplayName = (user) =>
  `${user.first_name || ""} ${user.last_name || ""}`.trim() ||
  user.nickname ||
  "Unnamed User";

const UserProfile = () => {
  const { userId } = useParams();
  const navigate = useNavigate();
  const { currentUser } = useAuth();
  const [user, setUser] = useState(null);
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [followerCount, setFollowerCount] = useState(0);
  const [followingCount, setFollowingCount] = useState(0);
  const [followStatus, setFollowStatus] = useState("unfollowed");

  useEffect(() => {
    // If trying to view own profile, redirect to /profile
    if (currentUser && userId === currentUser.id) {
      navigate("/profile");
      return;
    }

    const loadUserProfile = async () => {
      setLoading(true);
      setError("");

      try {
        const [userData, postsData, followersData, followingData] =
          await Promise.all([
            apiFetch(`/api/users/${userId}`),
            apiFetch(`/api/users/${userId}/posts`),
            apiFetch(`/api/followers/followers?user_id=${userId}`),
            apiFetch(`/api/followers/following?user_id=${userId}`),
          ]);

        setUser(userData);
        setPosts(Array.isArray(postsData) ? postsData : postsData?.data || []);
        setFollowerCount(
          Array.isArray(followersData)
            ? followersData.length
            : followersData?.data?.length || 0
        );
        setFollowingCount(
          Array.isArray(followingData)
            ? followingData.length
            : followingData?.data?.length || 0
        );

        // Determine follow status based on current user's relationship
        // This would ideally come from the API, but we can infer from the user data
        if (userData?.is_following) {
          setFollowStatus("following");
        } else if (userData?.follow_request_pending) {
          setFollowStatus("requested");
        } else {
          setFollowStatus("unfollowed");
        }
      } catch (err) {
        setError(
          err instanceof ApiError ? err.message : "Failed to load user profile"
        );
      } finally {
        setLoading(false);
      }
    };

    if (userId) {
      loadUserProfile();
    }
  }, [userId, currentUser, navigate]);

  const handleFollowStatusChange = (newStatus) => {
    setFollowStatus(newStatus);
    // Update follower count
    if (newStatus === "following") {
      setFollowerCount((prev) => prev + 1);
    } else if (newStatus === "unfollowed") {
      setFollowerCount((prev) => Math.max(0, prev - 1));
    }
  };

  if (loading) {
    return (
      <div className="profile-page">
        <div className="profile-page__inner">
          <div className="profile-skeleton profile-skeleton--header" />
          <div className="profile-skeleton profile-skeleton--row" />
          <div className="profile-skeleton profile-skeleton--row" />
        </div>
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="profile-page">
        <div className="profile-page__inner">
          <div className="profile-state profile-state--error">
            {error || "User not found"}
          </div>
        </div>
      </div>
    );
  }

  const displayName = toDisplayName(user);

  return (
    <div className="profile-page">
      <div className="profile-page__inner">
        <div className="profile-header">
          <div className="profile-header__cover" />

          <div className="profile-header__body">
            <div className="profile-header__avatar-wrap">
              <img
                src={user.avatar || avatarFallback}
                alt={`${displayName}'s avatar`}
                className="profile-header__avatar"
              />
            </div>

            <div className="profile-header__meta">
              <div className="profile-header__names">
                <h1>{displayName}</h1>
                {user.nickname && (
                  <p className="profile-header__nickname">@{user.nickname}</p>
                )}
                <div className="profile-header__badges">
                  <span
                    className={`profile-badge ${user.is_public ? "" : "is-private"}`}
                  >
                    {user.is_public ? "Public" : "Private"}
                  </span>
                </div>
              </div>

              <div className="profile-header__actions">
                <FollowAction
                  targetUserId={user.id}
                  initialStatus={followStatus}
                  isPrivate={!user.is_public}
                  onStatusChange={handleFollowStatusChange}
                />
              </div>
            </div>
          </div>

          <p
            className={`profile-header__bio ${user.about_me ? "" : "profile-header__bio--empty"}`}
          >
            {user.about_me || "This user hasn't written a bio yet."}
          </p>
        </div>

        <div className="profile-stats-card">
          <div className="profile-stats">
            <div className="profile-stats__item">
              <span className="profile-stats__count">{posts.length}</span>
              <span className="profile-stats__label">Posts</span>
            </div>
            <button
              type="button"
              className="profile-stats__item"
              onClick={() => {
                // Open followers modal
              }}
            >
              <span className="profile-stats__count">{followerCount}</span>
              <span className="profile-stats__label">Followers</span>
            </button>
            <button
              type="button"
              className="profile-stats__item"
              onClick={() => {
                // Open following modal
              }}
            >
              <span className="profile-stats__count">{followingCount}</span>
              <span className="profile-stats__label">Following</span>
            </button>
          </div>
        </div>

        <div className="profile-posts">
          {posts.length === 0 ? (
            <div className="profile-state">
              This user hasn't posted anything yet.
            </div>
          ) : (
            posts.map((post) => <Post key={post.id} post={post} />)
          )}
        </div>
      </div>
    </div>
  );
};

export default UserProfile;
```
11. Updated src/App.jsx
```jsx
import { Route, Routes } from "react-router";
import Layout from "./components/layout/Layout";
import Events from "./pages/Events.jsx";
import Groups from "./pages/Groups.jsx";
import Home from "./pages/Home.jsx";
import Friends from "./pages/Friends.jsx";
import Profile from "./pages/Profile.jsx";
import UserProfile from "./pages/UserProfile.jsx";
import Messages from "./pages/Messages.jsx";
import Notifications from "./pages/Notifications.jsx";
import PostDetail from "./pages/PostDetail.jsx";
import RegisterForm from "./components/RegisterForm";
import LoginForm from "./components/LoginForm";
import ProtectedRoute from "./components/ProtectedRoute";
import { useAuth } from "./context/useAuth.js";

function App() {
  const { isAuthenticated } = useAuth();
  return (
    <Routes>
      {/* If not authenticated, display this page without <Layout/>*/}
      {!isAuthenticated && <Route path="/post/:id" element={<PostDetail />} />}
      <Route path="/register" element={<RegisterForm />} />
      <Route path="/login" element={<LoginForm />} />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <Layout />
          </ProtectedRoute>
        }
      >
        <Route index element={<Home />} />
        <Route path="/post/:id" element={<PostDetail />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/user/:userId" element={<UserProfile />} />
        <Route path="/friends" element={<Friends />} />
        <Route path="/events" element={<Events />} />
        <Route path="/groups" element={<Groups />} />
        <Route path="/messages" element={<Messages />} />
        <Route path="/notifications" element={<Notifications />} />
      </Route>
    </Routes>
  );
}


export default App;
```
12. Updated src/components/Post.jsx (to make author clickable)
Add this to the existing Post component - replace the author avatar and name section:

```jsx
// Add this import at the top
import { useNavigate } from "react-router";

// In the Post component, find the post-header div and update it:
<div className="post-header">
  <img
    src={post?.author?.avatar ? post.author.avatar : avatar}
    alt="avatar"
    className="profile-photo"
    onClick={(e) => {
      e.stopPropagation();
      if (post?.author?.id) {
        navigate(`/user/${post.author.id}`);
      }
    }}
    style={{ cursor: "pointer" }}
  />

  <div className="post-bio">
    <h5 
      onClick={(e) => {
        e.stopPropagation();
        if (post?.author?.id) {
          navigate(`/user/${post.author.id}`);
        }
      }}
      style={{ cursor: "pointer" }}
    >
      {authorName}
    </h5>
    <small>{DateFormatter(post?.created_at, renderedAt)}</small>
  </div>
</div>
13. Add to src/styles/friends.css
Add this for the badge:

css
.friends-badge {
  display: inline-block;
  background: var(--hyperpop-violet);
  color: white;
  font-size: 0.7rem;
  padding: 2px 8px;
  border-radius: 12px;
  margin-left: 8px;
}

.friends-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
```
Now you have all the components needed for the follow system! The key integrations are:

FollowAction - The button component that handles follow/unfollow/cancel states

UserCard - Reusable card for displaying users with actions

FollowRequestsList - Displays pending follow requests with accept/reject

UserProfile - View other users' profiles with follow button

Friends - Updated to use the new components

Post - Updated to make author names/avatars clickable

App - Added /user/:userId route



Based on the dependencies between files, here's the recommended order to create the files:

Phase 1: Foundation (No Dependencies)
Create these first since they don't depend on any other new components:

`src/styles/user-card.css` - Base styles for user cards

`src/styles/follow-action.css` - Styles for follow buttons

`src/styles/follow-request-card.css` - Styles for request cards

`src/styles/follow-requests-list.css` - Styles for the list container

Phase 2: Core Components (Depend only on styles)
These components only need their CSS files:

`src/components/user/UserCard.jsx` - Uses user-card.css

`src/components/follow/FollowAction.jsx` - Uses follow-action.css

`src/components/follow/FollowRequestCard.jsx` - Uses follow-request-card.css

`src/components/follow/FollowRequestsList.jsx` - Uses follow-requests-list.css, imports FollowRequestCard

Phase 3: Page Updates (Depend on Phase 2 components)
These files import and use the new components:

`src/pages/Friends.jsx` - Replace with updated version that uses UserCard, FollowAction, and FollowRequestsList

`src/pages/UserProfile.jsx`- New file that uses UserCard (indirectly via FollowAction)

Phase 4: Routing and Integration (Depend on Phase 3)
`src/App.jsx` - Add the /user/:userId route

`src/components/Post.jsx` - Add click handlers for author name/avatar

Phase 5: Styles Addition (Optional but nice)
`src/styles/friends.css` - Add the .friends-badge style

Quick Setup Commands
Run these commands to create all the file placeholders at once:

```bash
# Create directories
mkdir -p src/components/user
mkdir -p src/components/follow

# Create all files (empty)
touch src/styles/user-card.css
touch src/styles/follow-action.css
touch src/styles/follow-request-card.css
touch src/styles/follow-requests-list.css
touch src/components/user/UserCard.jsx
touch src/components/follow/FollowAction.jsx
touch src/components/follow/FollowRequestCard.jsx
touch src/components/follow/FollowRequestsList.jsx
touch src/pages/UserProfile.jsx
```
Dependency Diagram
```text
Phase 1: CSS Files
    ↓
Phase 2: Components (UserCard, FollowAction, FollowRequestCard, FollowRequestsList)
    ↓
Phase 3: Pages (Friends, UserProfile)
    ↓
Phase 4: Routing (App, Post)
```
Testing Order
After creating files in this order, you can test incrementally:

After Phase 2: Test each component in isolation (use Storybook or a test page)

After Phase 3: Test the Friends page loads and shows requests/search results

After Phase 4: Test clicking a post author navigates to UserProfile, and follow/unfollow works

This order ensures you always have a working app and can test each piece as you go.