-- ==============================================================================
-- PostgreSQL Database Schema (DBPlay)
-- ==============================================================================

-- ------------------------------------------------------------------------------
-- Posts Table
-- ------------------------------------------------------------------------------
CREATE TABLE posts (
    post_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- In a real system, we'd have a user_id here as a foreign key
    -- user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    caption TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for efficient chronological timeline queries
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);

-- ------------------------------------------------------------------------------
-- Comments Table
-- ------------------------------------------------------------------------------
CREATE TABLE comments (
    comment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
    -- user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index to quickly load all comments for a specific post
CREATE INDEX idx_comments_post_id ON comments(post_id, created_at DESC);
