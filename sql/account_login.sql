CREATE TABLE IF NOT EXISTS `account_login_%s` (
    account VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    password_hash VARCHAR(100),
    version BIGINT NOT NULL DEFAULT 0
);
