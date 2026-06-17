CREATE TABLE groups (
    id TEXT PRIMARY KEY,
    creator_id TEXT,
    title TEXT NOT NULL,
    description TEXT,
    created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),

    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_groups_creator_created ON groups(creator_id, created_at);
CREATE INDEX idx_groups_created_at ON groups(created_at);

CREATE TABLE group_members (
    group_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending_invite', 'pending_request', 'accepted')),

    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_group_members_user_status ON group_members(user_id, status);
CREATE INDEX idx_group_members_group_status ON group_members(group_id, status);
