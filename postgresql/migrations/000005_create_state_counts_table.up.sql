CREATE TABLE IF NOT EXISTS state_counts (
    state_count_id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL,
    count INT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_state_counts_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
