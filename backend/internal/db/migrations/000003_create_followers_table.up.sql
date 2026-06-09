CREATE TABLE followers (
    follower_id TEXT NOT NULL,
    followee_id TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'accepted')),
    created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    
    PRIMARY KEY (follower_id, followee_id),
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followee_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (follower_id <> followee_id) -- Prevents a user from following themselves
);

CREATE INDEX idx_followers_followee_id ON followers(followee_id);