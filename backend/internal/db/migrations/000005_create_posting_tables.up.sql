CREATE TABLE posts (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    group_id TEXT,
    content TEXT,
    image_url TEXT,
    privacy TEXT NOT NULL CHECK (privacy IN ('public', 'almost_private', 'private')),
    comment_count INTEGER NOT NULL DEFAULT 0,
    like_count INTEGER NOT NULL DEFAULT 0,
    dislike_count INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    updated_at TEXT,
    deleted_at TEXT,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE INDEX idx_posts_created_at ON posts(created_at);
CREATE INDEX idx_posts_user_created ON posts(user_id, created_at);
CREATE INDEX idx_posts_group_created ON posts(group_id, created_at);
CREATE INDEX idx_posts_privacy_created ON posts(privacy, created_at);

CREATE TABLE post_audiences (
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,

    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_post_audiences_user ON post_audiences(user_id);

CREATE TABLE comments (
    id TEXT PRIMARY KEY,
    post_id TEXT NOT NULL,
    user_id TEXT,
    parent_comment_id TEXT,
    content TEXT,
    image_url TEXT,
    like_count INTEGER NOT NULL DEFAULT 0,
    dislike_count INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    deleted_at TEXT,
    updated_at TEXT,

    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (parent_comment_id) REFERENCES comments(id)
);

CREATE INDEX idx_comments_post_created ON comments(post_id, created_at);
CREATE INDEX idx_comments_parent_created ON comments(parent_comment_id, created_at);
CREATE INDEX idx_comments_user_created ON comments(user_id, created_at);

CREATE TABLE post_votes (
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    vote TEXT NOT NULL CHECK (vote IN ('like', 'dislike')),
    created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    updated_at TEXT,

    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_post_votes_user ON post_votes(user_id);
CREATE INDEX idx_post_votes_vote ON post_votes(vote);

CREATE TABLE comment_votes (
    comment_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    vote TEXT NOT NULL CHECK (vote IN ('like', 'dislike')),
    created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    updated_at TEXT,

    PRIMARY KEY (comment_id, user_id),
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_comment_votes_user ON comment_votes(user_id);
CREATE INDEX idx_comment_votes_vote ON comment_votes(vote);
