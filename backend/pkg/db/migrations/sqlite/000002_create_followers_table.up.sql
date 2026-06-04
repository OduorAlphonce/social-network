CREATE TABLE IF NOT EXISTS followers (
    follower_id TEXT NOT NULL,
    followee_id TEXT NOT NULL,
    status TEXT NOT NULL, -- 'pending' or 'accepted'
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, followee_id),
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followee_id) REFERENCES users(id) ON DELETE CASCADE
);
