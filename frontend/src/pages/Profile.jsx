import { useEffect, useState } from "react";
import { useAuth } from "../context/useAuth";
import { apiFetch, ApiError } from "../utils/api";
import ProfileHeader from "../components/profile/ProfileHeader";
import ProfileStats from "../components/profile/ProfileStats";
import ProfileTabs from "../components/profile/ProfileTabs";
import ProfileAbout from "../components/profile/ProfileAbout";
import ProfilePosts from "../components/profile/ProfilePosts";
import FollowListModal from "../components/profile/FollowList";
import ProfileUpdateForm from "../components/ProfileUpdateForm";
import "../styles/profile.css";

const TABS = [
  { id: "posts", label: "Posts" },
  { id: "about", label: "About" },
];

/**
 * Own-profile page: header (avatar/cover/bio/edit), stats, tabs for
 * posts and about info, plus modals for editing the profile and for
 * viewing the followers/following lists.
 */
const Profile = () => {
  const { currentUser, refresh } = useAuth();
  const [activeTab, setActiveTab] = useState("posts");
  const [isEditing, setIsEditing] = useState(false);
  const [followModal, setFollowModal] = useState(null); // "followers" | "following" | null
  const [followerCount, setFollowerCount] = useState(null);
  const [followingCount, setFollowingCount] = useState(null);
  const [postCount, setPostCount] = useState(0);
  const [avatarError, setAvatarError] = useState("");

  const userId = currentUser?.id;

  const loadCounts = async () => {
    if (!userId) return;
    try {
      const [followers, following] = await Promise.all([
        apiFetch(`/api/followers/followers?user_id=${userId}`),
        apiFetch(`/api/followers/following?user_id=${userId}`),
      ]);
      setFollowerCount(
        Array.isArray(followers) ? followers.length : (followers?.data || []).length
      );
      setFollowingCount(
        Array.isArray(following) ? following.length : (following?.data || []).length
      );
    } catch {
      // Counts are a nice-to-have; failing silently keeps the page usable.
      setFollowerCount(0);
      setFollowingCount(0);
    }
  };

  const loadPostCount = async () => {
    if (!userId) return;
    try {
      const result = await apiFetch(`/api/users/${userId}/posts`);
      const list = Array.isArray(result) ? result : result?.data || [];
      setPostCount(list.length);
    } catch {
      setPostCount(0);
    }
  };

  // Load counts once we know who the current user is.
  useEffect(() => {
    loadCounts();
    loadPostCount();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [userId]);

  const handleAvatarChange = async (file) => {
    setAvatarError("");
    const data = new FormData();
    data.append("avatar", file);

    try {
      await apiFetch("/api/users/update", { method: "PATCH", body: data });
      await refresh();
    } catch (err) {
      setAvatarError(
        err instanceof ApiError ? err.message : "Unable to update avatar."
      );
    }
  };

  if (!currentUser) {
    return (
      <div className="profile-page">
        <div className="profile-page__inner">
          <div className="profile-skeleton profile-skeleton--header" />
        </div>
      </div>
    );
  }

  return (
    <div className="profile-page">
      <div className="profile-page__inner">
        <ProfileHeader
          user={currentUser}
          isOwnProfile
          onEdit={() => setIsEditing(true)}
          onAvatarChange={handleAvatarChange}
        />

        {avatarError && (
          <div className="profile-state profile-state--error">
            {avatarError}
          </div>
        )}

        <div className="profile-stats-card">
          <ProfileStats
            postCount={postCount}
            followerCount={followerCount ?? 0}
            followingCount={followingCount ?? 0}
            onShowFollowers={() => setFollowModal("followers")}
            onShowFollowing={() => setFollowModal("following")}
          />
        </div>

        <ProfileTabs
          tabs={TABS}
          activeTab={activeTab}
          onChange={setActiveTab}
        />

        {activeTab === "posts" ? (
          <ProfilePosts userId={userId} />
        ) : (
          <ProfileAbout user={currentUser} isOwnProfile />
        )}
      </div>

      {isEditing && (
        <div className="profile-edit-overlay">
          <ProfileUpdateForm onClose={() => setIsEditing(false)} />
        </div>
      )}

      {followModal && (
        <FollowListModal
          userId={userId}
          type={followModal}
          onClose={() => setFollowModal(null)}
        />
      )}
    </div>
  );
};

export default Profile;