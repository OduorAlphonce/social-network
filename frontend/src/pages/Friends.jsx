import { useEffect, useState } from "react";
import { useAuth } from "../context/useAuth";
import { apiFetch, ApiError } from "../utils/api";
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
          <h2>Pending requests</h2>
        </div>

        {loading ? (
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
        )}
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
