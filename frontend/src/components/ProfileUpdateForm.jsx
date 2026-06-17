import { useState } from "react";
import { useAuth } from "../context/useAuth";
import "../styles/RegisterForm.css"; // Reusing styles for consistency
import { apiFetch } from "../utils/api";

const getProfileFormData = (currentUser) => ({
  email: currentUser.email || "",
  current_password: "",
  new_password: "",
  first_name: currentUser.first_name || "",
  last_name: currentUser.last_name || "",
  date_of_birth: currentUser.date_of_birth || "",
  nickname: currentUser.nickname || "",
  about_me: currentUser.about_me || "",
  is_public: currentUser.is_public || false,
});

const ProfileUpdateFields = ({ currentUser, refresh }) => {
  const [formData, setFormData] = useState(() =>
    getProfileFormData(currentUser)
  );
  const [avatar, setAvatar] = useState(null);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: type === "checkbox" ? checked : value,
    }));
  };

  const handleFileChange = (e) => {
    setAvatar(e.target.files[0]);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    const data = new FormData();
    Object.keys(formData).forEach((key) => {
      if (formData[key] !== "" || key === "is_public") {
        data.append(key, formData[key]);
      }
    });
    if (avatar) {
      data.append("avatar", avatar);
    }

    try {
      await apiFetch("/api/users/update", {
        method: "PATCH",
        body: data,
      });

      setSuccess("Profile updated successfully!");
      await refresh();
      setFormData((prev) => ({
        ...prev,
        current_password: "",
        new_password: "",
      }));
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="register-form-container">
      <h2>Update Profile</h2>
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}
      <form onSubmit={handleSubmit} className="register-form">
        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
          />
        </div>

        <div className="form-row">
          <div className="form-group">
            <label htmlFor="first_name">First Name</label>
            <input
              type="text"
              id="first_name"
              name="first_name"
              value={formData.first_name}
              onChange={handleChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="last_name">Last Name</label>
            <input
              type="text"
              id="last_name"
              name="last_name"
              value={formData.last_name}
              onChange={handleChange}
            />
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="date_of_birth">Date of Birth</label>
          <input
            type="date"
            id="date_of_birth"
            name="date_of_birth"
            value={formData.date_of_birth}
            onChange={handleChange}
          />
        </div>

        <div className="form-group">
          <label htmlFor="nickname">Nickname</label>
          <input
            type="text"
            id="nickname"
            name="nickname"
            value={formData.nickname}
            onChange={handleChange}
          />
        </div>

        <div className="form-group">
          <label htmlFor="about_me">About Me</label>
          <textarea
            id="about_me"
            name="about_me"
            value={formData.about_me}
            onChange={handleChange}
          />
        </div>

        <div className="form-group checkbox-group">
          <label>
            <input
              type="checkbox"
              name="is_public"
              checked={formData.is_public}
              onChange={handleChange}
            />
            Public Profile
          </label>
        </div>

        <div className="form-group">
          <label htmlFor="avatar">Avatar</label>
          <input
            type="file"
            id="avatar"
            name="avatar"
            accept="image/*"
            onChange={handleFileChange}
          />
        </div>

        <hr />
        <p className="form-hint">
          Fill these only if you want to change email or password
        </p>

        <div className="form-group">
          <label htmlFor="current_password">
            Current Password (Required for sensitive changes)
          </label>
          <input
            type="password"
            id="current_password"
            name="current_password"
            value={formData.current_password}
            onChange={handleChange}
          />
        </div>

        <div className="form-group">
          <label htmlFor="new_password">New Password</label>
          <input
            type="password"
            id="new_password"
            name="new_password"
            value={formData.new_password}
            onChange={handleChange}
          />
        </div>

        <button type="submit" className="submit-btn">
          Update Profile
        </button>
      </form>
    </div>
  );
};

const ProfileUpdateForm = () => {
  const { currentUser, refresh } = useAuth();

  if (!currentUser) {
    return <div>Loading...</div>;
  }

  return (
    <ProfileUpdateFields
      key={currentUser.id || currentUser.email}
      currentUser={currentUser}
      refresh={refresh}
    />
  );
};

export default ProfileUpdateForm;
