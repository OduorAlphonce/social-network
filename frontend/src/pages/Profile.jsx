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

// MOCK DATA - Realistic user profile
const MOCK_USER = {
  id: "user_12345",
  username: "alex_morgan",
  displayName: "Alex Morgan",
  email: "alex@example.com",
  bio: "Full-stack developer | Coffee enthusiast | 📍 NYC",
  avatar: "https://i.pravatar.cc/150?u=alex_morgan",
  coverImage: "https://picsum.photos/1200/400?random=1",
  location: "New York, NY",
  website: "https://alexmorgan.dev",
  joinedDate: "2023-08-15T14:30:00Z",
  socialLinks: {
    twitter: "@alexmorgan",
    github: "alexmorgan",
    linkedin: "alexmorgan"
  },
  interests: ["React", "TypeScript", "GraphQL", "Docker", "Python"]
};

// MOCK DATA - Follower/Following lists
const MOCK_FOLLOWERS = [
  { id: "user_001", username: "sarah_dev", displayName: "Sarah Chen", avatar: "https://i.pravatar.cc/150?u=sarah" },
  { id: "user_002", username: "mike_codes", displayName: "Mike Johnson", avatar: "https://i.pravatar.cc/150?u=mike" },
  { id: "user_003", username: "emily_tech", displayName: "Emily Rodriguez", avatar: "https://i.pravatar.cc/150?u=emily" },
  { id: "user_004", username: "david_ai", displayName: "David Kim", avatar: "https://i.pravatar.cc/150?u=david" },
  { id: "user_005", username: "lisa_py", displayName: "Lisa Patel", avatar: "https://i.pravatar.cc/150?u=lisa" },
];

const MOCK_FOLLOWING = [
  { id: "user_006", username: "react_guru", displayName: "React Master", avatar: "https://i.pravatar.cc/150?u=react" },
  { id: "user_007", username: "design_queen", displayName: "Design Queen", avatar: "https://i.pravatar.cc/150?u=design" },
  { id: "user_008", username: "cloud_ninja", displayName: "Cloud Ninja", avatar: "https://i.pravatar.cc/150?u=cloud" },
  { id: "user_009", username: "js_wizard", displayName: "JS Wizard", avatar: "https://i.pravatar.cc/150?u=js" },
  { id: "user_010", username: "rustacean", displayName: "Rustacean", avatar: "https://i.pravatar.cc/150?u=rust" },
  { id: "user_011", username: "devops_gal", displayName: "DevOps Gal", avatar: "https://i.pravatar.cc/150?u=devops" },
];

// MOCK DATA - Posts
const MOCK_POSTS = [
  {
    id: "post_001",
    content: "Just launched my new portfolio website! Built with React and Tailwind CSS. Check it out at https://alexmorgan.dev 🚀",
    createdAt: "2026-06-19T10:30:00Z",
    likes: 42,
    comments: 8,
    shares: 3,
    images: ["https://picsum.photos/800/600?random=10"]
  },
  {
    id: "post_002",
    content: "Finally solved the N+1 query problem in GraphQL. The key was using DataLoader properly! #graphql #performance",
    createdAt: "2026-06-18T16:45:00Z",
    likes: 28,
    comments: 5,
    shares: 2,
    images: []
  },
  {
    id: "post_003",
    content: "Beautiful day for coding in the park ☀️💻",
    createdAt: "2026-06-17T09:15:00Z",
    likes: 15,
    comments: 3,
    shares: 1,
    images: ["https://picsum.photos/800/600?random=20"]
  },
  {
    id: "post_004",
    content: "Just published a new blog post: 'Advanced TypeScript Patterns for React Developers'. Would love your feedback! 🧵",
    createdAt: "2026-06-15T12:00:00Z",
    likes: 56,
    comments: 12,
    shares: 7,
    images: []
  }
];

// MOCK DATA - Counts
const MOCK_COUNTS = {
  followers: 1247,
  following: 342,
  posts: 48
};

// Mock implementations for the component's functions
const mockRefresh = () => Promise.resolve();

const Profile = () => {
  // const { currentUser, refresh } = useAuth();
  const [activeTab, setActiveTab] = useState("posts");
  const [isEditing, setIsEditing] = useState(false);
  const [followModal, setFollowModal] = useState(null);
  const [followerCount, setFollowerCount] = useState(MOCK_COUNTS.followers);
  const [followingCount, setFollowingCount] = useState(MOCK_COUNTS.following);
  const [postCount, setPostCount] = useState(MOCK_COUNTS.posts);
  const [avatarError, setAvatarError] = useState("");

  // Use mock user data when currentUser is null
  const user = MOCK_USER; // || currentUser
  const userId = user.id;

  // Mock loading functions
  const loadCounts = async () => {
    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 300));
    setFollowerCount(MOCK_COUNTS.followers);
    setFollowingCount(MOCK_COUNTS.following);
  };

  const loadPostCount = async () => {
    await new Promise(resolve => setTimeout(resolve, 200));
    setPostCount(MOCK_COUNTS.posts);
  };

  useEffect(() => {
    loadCounts();
    loadPostCount();
  }, [userId]);

  const handleAvatarChange = async (file) => {
    setAvatarError("");
    const data = new FormData();
    data.append("avatar", file);

    try {
      await new Promise(resolve => setTimeout(resolve, 500));
      // Simulate successful avatar update
      await mockRefresh();
    } catch (err) {
      setAvatarError(
        err instanceof ApiError ? err.message : "Unable to update avatar."
      );
    }
  };

  // Mock version of ProfilePosts that uses mock data
  const MockProfilePosts = ({ userId }) => {
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
      const loadPosts = async () => {
        await new Promise(resolve => setTimeout(resolve, 400));
        setPosts(MOCK_POSTS);
        setLoading(false);
      };
      loadPosts();
    }, [userId]);

    return (
      <div className="profile-posts">
        {loading ? (
          <div className="profile-skeleton profile-skeleton--posts">
            {[1, 2, 3].map((i) => (
              <div key={i} className="post-skeleton" />
            ))}
          </div>
        ) : (
          <>
            {posts.length === 0 ? (
              <div className="profile-state profile-state--empty">
                <p>No posts yet. Start sharing your thoughts!</p>
              </div>
            ) : (
              posts.map((post) => (
                <div key={post.id} className="post-card">
                  <div className="post-card__header">
                    <img 
                      src={user.avatar} 
                      alt={user.displayName}
                      className="post-card__avatar"
                    />
                    <div className="post-card__meta">
                      <span className="post-card__name">{user.displayName}</span>
                      <span className="post-card__username">@{user.username}</span>
                      <span className="post-card__time">
                        {new Date(post.createdAt).toLocaleDateString('en-US', {
                          month: 'short',
                          day: 'numeric'
                        })}
                      </span>
                    </div>
                  </div>
                  <div className="post-card__content">{post.content}</div>
                  {post.images.length > 0 && (
                    <div className="post-card__images">
                      {post.images.map((img, idx) => (
                        <img key={idx} src={img} alt={`Post image ${idx + 1}`} />
                      ))}
                    </div>
                  )}
                  <div className="post-card__actions">
                    <button>❤️ {post.likes}</button>
                    <button>💬 {post.comments}</button>
                    <button>🔁 {post.shares}</button>
                  </div>
                </div>
              ))
            )}
          </>
        )}
      </div>
    );
  };

  // Mock version of ProfileAbout that uses mock data
  const MockProfileAbout = ({ user }) => (
    <div className="profile-about">
      <div className="profile-about__bio">
        <h3>Bio</h3>
        <p>{user.bio}</p>
      </div>
      
      <div className="profile-about__details">
        <h3>Details</h3>
        <div className="profile-about__detail-item">
          <span className="icon">📍</span>
          <span>{user.location}</span>
        </div>
        <div className="profile-about__detail-item">
          <span className="icon">🌐</span>
          <a href={user.website} target="_blank" rel="noopener noreferrer">
            {user.website}
          </a>
        </div>
        <div className="profile-about__detail-item">
          <span className="icon">📅</span>
          <span>Joined {new Date(user.joinedDate).toLocaleDateString('en-US', {
            month: 'long',
            year: 'numeric'
          })}</span>
        </div>
      </div>

      <div className="profile-about__interests">
        <h3>Interests</h3>
        <div className="profile-about__tags">
          {user.interests.map((interest) => (
            <span key={interest} className="interest-tag">
              #{interest}
            </span>
          ))}
        </div>
      </div>

      <div className="profile-about__social">
        <h3>Social Links</h3>
        {user.socialLinks.twitter && (
          <div className="profile-about__social-link">
            <span className="icon">🐦</span>
            <span>{user.socialLinks.twitter}</span>
          </div>
        )}
        {user.socialLinks.github && (
          <div className="profile-about__social-link">
            <span className="icon">💻</span>
            <span>{user.socialLinks.github}</span>
          </div>
        )}
        {user.socialLinks.linkedin && (
          <div className="profile-about__social-link">
            <span className="icon">🔗</span>
            <span>{user.socialLinks.linkedin}</span>
          </div>
        )}
      </div>
    </div>
  );

  return (
    <div className="profile-page">
      <div className="profile-page__inner">
        <ProfileHeader
          user={user}
          // isOwnProfile={!!currentUser}
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
          <MockProfilePosts userId={userId} />
        ) : (
          <MockProfileAbout user={user} />
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
          followersData={MOCK_FOLLOWERS}
          followingData={MOCK_FOLLOWING}
          onClose={() => setFollowModal(null)}
        />
      )}
    </div>
  );
};

export default Profile;