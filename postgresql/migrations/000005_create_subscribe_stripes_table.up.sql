CREATE TABLE IF NOT EXISTS subscribe_stripes (
    subscribe_stripe_id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL,
    subscribe_id VARCHAR(26) NOT NULL,
    session_id VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_subscribe_stripes_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_subscribe_stripes_subscribe_id
        FOREIGN KEY (subscribe_id)
        REFERENCES subscribes(subscribe_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
