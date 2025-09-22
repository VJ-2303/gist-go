CREATE TABLE IF NOT EXISTS share_links (
    id bigserial PRIMARY KEY,

    post_id bigint NOT NULL,
    CONSTRAINT fk_post
        FOREIGN KEY(post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE,

    token text UNIQUE NOT NULL,

    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    expires_at timestamp(0) with time zone
);
