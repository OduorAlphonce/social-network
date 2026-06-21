import { useState } from "react";
import "../styles/RegisterForm.css";
import { useNavigate } from "react-router";
import { apiFetch } from "../utils/api";

const INITIAL_FORM = {
  email: "",
  password: "",
  first_name: "",
  last_name: "",
  date_of_birth: "",
  nickname: "",
  about_me: "",
  is_public: false,
};

const RegisterForm = () => {
  const [step, setStep] = useState(1);
  const [formData, setFormData] = useState(INITIAL_FORM);
  const [avatar, setAvatar] = useState(null);
  const [avatarPreview, setAvatarPreview] = useState(null);
  const [dragActive, setDragActive] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const navigate = useNavigate();  

  // ---------- Handlers ----------

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: type === "checkbox" ? checked : value,
    }));
  };

  const handleFileChange = (e) => {
    const file = e.target.files?.[0];
    processFile(file);
  };

  const processFile = (file) => {
    if (!file) return;

    if (!file.type.startsWith("image/")) {
      setError("Please upload an image file.");
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      setError("Image must be smaller than 5MB.");
      return;
    }

    setError("");
    setAvatar(file);
    setAvatarPreview(URL.createObjectURL(file));
  };

  // Drag-and-drop handlers
  const handleDrag = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    const file = e.dataTransfer.files?.[0];
    processFile(file);
  };

  const handleRemoveAvatar = () => {
    setAvatar(null);
    setAvatarPreview(null);
  };

  // ---------- Validation ----------

  const validateStep1 = () => {
    if (!formData.email.trim()) return "Email is required.";
    if (!formData.password.trim()) return "Password is required.";
    if (formData.password.length < 8)
      return "Password must be at least 8 characters.";
    if (!formData.first_name.trim()) return "First name is required.";
    if (!formData.last_name.trim()) return "Last name is required.";
    if (!formData.date_of_birth) return "Date of birth is required.";
    return null;
  };

  const handleNext = () => {
    const validationError = validateStep1();
    if (validationError) {
      setError(validationError);
      return;
    }
    setError("");
    setStep(2);
  };

  const handleBack = () => {
    setError("");
    setStep(1);
  };

  // ---------- Submit ----------

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");
    setIsSubmitting(true);

    const data = new FormData();
    Object.keys(formData).forEach((key) => {
      data.append(key, formData[key]);
    });
    if (avatar) {
      data.append("avatar", avatar);
    }

    try {
      await apiFetch("/api/users/register", {
        method: "POST",
        body: data,
      });

      // Redirect to login page with a success flag
      navigate("/login", { 
        state: { registered: true },
        replace: true 
      });
      
    } catch (err) {
      setError(err.message || "Registration failed. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };



  // ---------- Step 1: Required Details ----------

  if (step === 1) {
    return (
      <div className="register-form-container">
        <div className="register-steps">
          <span className="register-steps__dot is-active" />
          <span className="register-steps__line" />
          <span className="register-steps__dot" />
        </div>
        <h2>Create Account</h2>
        <p className="register-subtitle">Step 1 of 2 — The essentials</p>

        {error && <div className="error-message">{error}</div>}

        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleNext();
          }}
          className="register-form"
        >
          <div className="form-group">
            <label htmlFor="email">Email Address *</label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              required
              placeholder="you@example.com"
              autoComplete="email"
            />
          </div>

          <div className="form-group">
            <label htmlFor="password">Password *</label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              placeholder="Minimum 8 characters"
              autoComplete="new-password"
            />
          </div>

          <div className="form-row">
            <div className="form-group">
              <label htmlFor="first_name">First Name *</label>
              <input
                type="text"
                id="first_name"
                name="first_name"
                value={formData.first_name}
                onChange={handleChange}
                required
                placeholder="First name"
              />
            </div>
            <div className="form-group">
              <label htmlFor="last_name">Last Name *</label>
              <input
                type="text"
                id="last_name"
                name="last_name"
                value={formData.last_name}
                onChange={handleChange}
                required
                placeholder="Last name"
              />
            </div>
          </div>

          <div className="form-group">
            <label htmlFor="date_of_birth">Date of Birth *</label>
            <input
              type="date"
              id="date_of_birth"
              name="date_of_birth"
              value={formData.date_of_birth}
              onChange={handleChange}
              required
            />
          </div>

          <button type="submit" className="submit-btn">
            Continue
          </button>
        </form>

        <div className="register-footer">
          <p>
            Already have an account? <a href="/login">Log in</a>
          </p>
        </div>
      </div>
    );
  }

  // ---------- Step 2: Optional Details + Avatar ----------

  return (
    <div className="register-form-container">
      <div className="register-steps">
        <span className="register-steps__dot is-complete" />
        <span className="register-steps__line is-complete" />
        <span className="register-steps__dot is-active" />
      </div>
      <h2>Complete Your Profile</h2>
      <p className="register-subtitle">Step 2 of 2 — Optional details</p>

      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}

      <form onSubmit={handleSubmit} className="register-form">
        {/* Avatar Upload */}
        <div className="form-group">
          <label>Profile Picture</label>
          <div
            className={`avatar-upload-zone ${dragActive ? "avatar-upload-zone--active" : ""} ${avatarPreview ? "avatar-upload-zone--has-image" : ""}`}
            onDragEnter={handleDrag}
            onDragLeave={handleDrag}
            onDragOver={handleDrag}
            onDrop={handleDrop}
            onClick={() => !avatarPreview && document.getElementById("avatar").click()}
          >
            {avatarPreview ? (
              <div className="avatar-upload-zone__preview">
                <img src={avatarPreview} alt="Avatar preview" />
                <div className="avatar-upload-zone__overlay">
                  <span>Change</span>
                </div>
              </div>
            ) : (
              <div className="avatar-upload-zone__placeholder">
                <svg
                  width="48"
                  height="48"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="1.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h7" />
                  <line x1="16" y1="5" x2="22" y2="5" />
                  <line x1="19" y1="2" x2="19" y2="8" />
                  <circle cx="9" cy="9" r="2" />
                  <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21" />
                </svg>
                <span>Drag & drop your photo here</span>
                <span className="avatar-upload-zone__hint">
                  or click to browse — JPG, PNG, max 5MB
                </span>
              </div>
            )}
            <input
              type="file"
              id="avatar"
              name="avatar"
              accept="image/*"
              onChange={handleFileChange}
              className="avatar-upload-zone__input"
            />
          </div>
          {avatarPreview && (
            <button
              type="button"
              className="avatar-remove-btn"
              onClick={handleRemoveAvatar}
            >
              Remove photo
            </button>
          )}
        </div>

        <div className="form-group">
          <label htmlFor="nickname">Nickname</label>
          <input
            type="text"
            id="nickname"
            name="nickname"
            value={formData.nickname}
            onChange={handleChange}
            placeholder="How should people know you?"
          />
        </div>

        <div className="form-group">
          <label htmlFor="about_me">About Me</label>
          <textarea
            id="about_me"
            name="about_me"
            value={formData.about_me}
            onChange={handleChange}
            placeholder="Tell the world a little about yourself..."
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
            Make my profile public
          </label>
        </div>

        <div className="form-actions">
          <button
            type="button"
            className="submit-btn submit-btn--ghost"
            onClick={handleBack}
          >
            Back
          </button>
          <button
            type="submit"
            className="submit-btn"
            disabled={isSubmitting}
          >
            {isSubmitting ? "Creating Account..." : "Complete Registration"}
          </button>
        </div>
      </form>
    </div>
  );
};

export default RegisterForm;