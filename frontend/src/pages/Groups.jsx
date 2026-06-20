import { useEffect, useState } from "react";
import { apiFetch } from "../utils/api";
import "../styles/LoginForm.css"; // Reuse some card and form classes for consistency

const Groups = () => {
  const [groups, setGroups] = useState([]);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [showCreateModal, setShowCreateModal] = useState(false);

  const fetchGroups = async () => {
    try {
      const data = await apiFetch("/api/groups");
      setGroups(data || []);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchGroups();
  }, []);

  const handleCreateGroup = async (e) => {
    e.preventDefault();
    if (!title) return;
    setLoading(true);
    setError("");
    try {
      await apiFetch("/api/groups", {
        method: "POST",
        body: { title, description },
      });
      setTitle("");
      setDescription("");
      setShowCreateModal(false);
      fetchGroups();
    } catch (err) {
      setError(err.message || "Failed to create group");
    } finally {
      setLoading(false);
    }
  };

  const handleJoinGroup = async (gId) => {
    try {
      await apiFetch(`/api/groups/${gId}/join`, { method: "POST" });
      fetchGroups();
    } catch (err) {
      alert(err.message || "Failed to join group");
    }
  };

  return (
    <div style={{ padding: "20px", maxWidth: "800px", margin: "0 auto" }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "20px" }}>
        <h2 style={{ margin: 0 }}>Groups</h2>
        <button
          onClick={() => setShowCreateModal(true)}
          style={{
            background: "linear-gradient(135deg, #667eea, #764ba2)",
            color: "white",
            border: "none",
            padding: "10px 20px",
            borderRadius: "8px",
            fontWeight: "bold",
            cursor: "pointer",
            transition: "all 0.3s ease",
          }}
        >
          Create Group
        </button>
      </div>

      {showCreateModal && (
        <div style={{
          position: "fixed", top: 0, left: 0, right: 0, bottom: 0,
          backgroundColor: "rgba(0,0,0,0.5)", display: "flex", justifyContent: "center", alignItems: "center", zIndex: 1000
        }}>
          <div style={{ backgroundColor: "#1e1e1e", padding: "30px", borderRadius: "12px", width: "400px", color: "white" }}>
            <h3 style={{ marginTop: 0 }}>Create New Group</h3>
            {error && <div style={{ color: "#ff6b6b", marginBottom: "15px" }}>{error}</div>}
            <form onSubmit={handleCreateGroup}>
              <div style={{ marginBottom: "15px" }}>
                <label style={{ display: "block", marginBottom: "5px" }}>Title</label>
                <input
                  type="text"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  style={{ width: "100%", padding: "10px", borderRadius: "6px", border: "1px solid #444", backgroundColor: "#2e2e2e", color: "white" }}
                  required
                />
              </div>
              <div style={{ marginBottom: "20px" }}>
                <label style={{ display: "block", marginBottom: "5px" }}>Description</label>
                <textarea
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  style={{ width: "100%", padding: "10px", borderRadius: "6px", border: "1px solid #444", backgroundColor: "#2e2e2e", color: "white", height: "80px" }}
                />
              </div>
              <div style={{ display: "flex", justifyContent: "flex-end", gap: "10px" }}>
                <button
                  type="button"
                  onClick={() => setShowCreateModal(false)}
                  style={{ padding: "8px 16px", borderRadius: "6px", border: "1px solid #555", backgroundColor: "transparent", color: "white", cursor: "pointer" }}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  style={{
                    padding: "8px 16px", borderRadius: "6px", border: "none",
                    background: "linear-gradient(135deg, #667eea, #764ba2)", color: "white", cursor: "pointer"
                  }}
                >
                  {loading ? "Creating..." : "Create"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(240px, 1fr))", gap: "20px" }}>
        {groups.map((group) => (
          <div
            key={group.id}
            style={{
              backgroundColor: "#1e1e1e",
              borderRadius: "12px",
              padding: "20px",
              boxShadow: "0 4px 6px rgba(0,0,0,0.1)",
              border: "1px solid #2e2e2e",
              display: "flex",
              flexDirection: "column",
              justifyContent: "space-between"
            }}
          >
            <div>
              <h3 style={{ margin: "0 0 10px 0", color: "#667eea" }}>{group.title}</h3>
              <p style={{ color: "#aaa", fontSize: "0.9rem", margin: "0 0 15px 0" }}>{group.description || "No description provided."}</p>
            </div>

            {group.status === "accepted" ? (
              <span style={{ color: "#2ecc71", fontWeight: "bold", textAlign: "center", display: "block" }}>Member</span>
            ) : group.status === "pending_request" ? (
              <span style={{ color: "#f1c40f", fontWeight: "bold", textAlign: "center", display: "block" }}>Request Pending</span>
            ) : group.status === "pending_invite" ? (
              <div style={{ display: "flex", gap: "10px" }}>
                <button
                  onClick={() => handleJoinGroup(group.id)} // accepts invite
                  style={{
                    flex: 1, backgroundColor: "#2ecc71", color: "white", border: "none",
                    padding: "8px", borderRadius: "6px", cursor: "pointer", fontWeight: "bold"
                  }}
                >
                  Accept
                </button>
              </div>
            ) : (
              <button
                onClick={() => handleJoinGroup(group.id)}
                style={{
                  backgroundColor: "#667eea", color: "white", border: "none",
                  padding: "8px", borderRadius: "6px", cursor: "pointer", fontWeight: "bold"
                }}
              >
                Join Group
              </button>
            )}
          </div>
        ))}
        {groups.length === 0 && (
          <div style={{ color: "#888", gridColumn: "1/-1", textAlign: "center", padding: "40px" }}>
            No groups available. Create one to get started!
          </div>
        )}
      </div>
    </div>
  );
};

export default Groups;
