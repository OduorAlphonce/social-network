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

/**
 * Friends page - displays pending follow requests and discoverable users.
 * Uses modern components (FollowRequestsList, UserCard, FollowAction) with
 * fallback to legacy implementation for compatibility.
 */
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
  
  // Feature flag to toggle between new and legacy UI
  // Set to false to use legacy implementation, true to use new components
  const [useNewComponents] = useState(true);

  // Load initial data
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

  // Refresh pending requests
  const refreshPending = async () => {
    try {
      const pending = await apiFetch("/api/followers/pending");
      setPendingRequests(Array.isArray(pending) ? pending : pending || []);
      setRequestCount(pending?.length || 0);
    } catch {
      // don't block the user if pending refresh fails
    }
  };

  // Load/search suggestions
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

  // --- Legacy Handlers (kept for backwards compatibility) ---
  const handleSendRequest = async (userId) => {
    setActionMessage("");
    setError("");
    setSubmitting((prev) => ({ ...prev, [userId]: true }));

    try {
      await apiFetch("/api/followers/follow", {
        method: "POST",
        body: { following_id: userId },
      });
      setActionMessage("Friend request sent.");
      await loadSuggestions(query);
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Unable to send friend request."
      );
    } finally {
      setSubmitting((prev) => ({ ...prev, [userId]: false }));
    }
  };

  const handleAcceptRequest = async (followerId) => {
    setActionMessage("");
    setError("");
    setSubmitting((prev) => ({ ...prev, [followerId]: true }));

    try {
      await apiFetch("/api/followers/accept", {
        method: "POST",
        body: { follower_id: followerId },
      });
      setActionMessage("Follow request accepted.");
      await refreshPending();
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Unable to accept request."
      );
    } finally {
      setSubmitting((prev) => ({ ...prev, [followerId]: false }));
    }
  };

  const handleRejectRequest = async (followerId) => {
    setActionMessage("");
    setError("");
    setSubmitting((prev) => ({ ...prev, [followerId]: true }));

    try {
      await apiFetch("/api/followers/reject", {
        method: "POST",
        body: { follower_id: followerId },
      });
      setActionMessage("Follow request rejected.");
      await refreshPending();
    } catch (err) {
      setError(
        err instanceof ApiError ? err.message : "Unable to reject request."
      );
    } finally {
      setSubmitting((prev) => ({ ...prev, [followerId]: false }));
    }
  };

  // --- New Handler (for FollowAction component) ---
  const handleFollowStatusChange = (userId, newStatus) => {
    // Update the suggested users list to reflect new status
    setSuggestedUsers((prev) =>
      prev.map((user) => {
        if (user.id === userId) {
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

      {/* Pending Requests Section */}
      <section className="friends-section">
        <div className="friends-section__header">
          <h2>
            Pending requests{" "}
            {requestCount > 0 && (
              <span className="friends-badge">{requestCount}</span>
            )}
          </h2>
        </div>

        {useNewComponents ? (
          // NEW: Use FollowRequestsList component
          <FollowRequestsList onRequestCountChange={setRequestCount} />
        ) : (
          // LEGACY: Manual pending requests rendering
          loading ? (
            <div className="profile-skeleton profile-skeleton--row" />
          ) : pendingRequests.length === 0 ? (
            <div className="profile-state">No pending follow requests.</div>
          ) : (
            <ul className="friends-list">
              {pendingRequests.map((user) => (
                <li key={user.id} className="friends-card">
                  <img
                    src={user.avatar || avatarFallback}
                    alt={`${toDisplayName(user)}'s avatar`}
                    className="friends-card__avatar"
                  />
                  <div className="friends-card__content">
                    <p className="friends-card__name">{toDisplayName(user)}</p>
                    {user.nickname && (
                      <p className="friends-card__handle">@{user.nickname}</p>
                    )}
                  </div>
                  <div className="friends-card__actions">
                    <button
                      type="button"
                      className="profile-btn profile-btn--primary"
                      disabled={Boolean(submitting[user.id])}
                      onClick={() => handleAcceptRequest(user.id)}
                    >
                      {submitting[user.id] ? "Accepting..." : "Accept"}
                    </button>
                    <button
                      type="button"
                      className="profile-btn profile-btn--ghost"
                      disabled={Boolean(submitting[user.id])}
                      onClick={() => handleRejectRequest(user.id)}
                    >
                      Reject
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          )
        )}
      </section>

      {/* Discover Profiles Section */}
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
        ) : useNewComponents ? (
          // NEW: Use UserCard with FollowAction
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
        ) : (
          // LEGACY: Manual user list rendering
          <ul className="friends-list">
            {suggestedUsers.map((user) => (
              <li key={user.id} className="friends-card">
                <img
                  src={user.avatar || avatarFallback}
                  alt={`${toDisplayName(user)}'s avatar`}
                  className="friends-card__avatar"
                />
                <div className="friends-card__content">
                  <p className="friends-card__name">{toDisplayName(user)}</p>
                  {user.nickname && (
                    <p className="friends-card__handle">@{user.nickname}</p>
                  )}
                </div>
                <button
                  type="button"
                  className="profile-btn profile-btn--primary"
                  disabled={Boolean(submitting[user.id])}
                  onClick={() => handleSendRequest(user.id)}
                >
                  {submitting[user.id] ? "Sending..." : "Send request"}
                </button>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
};

export default Friends;