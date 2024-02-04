CREATE TABLE IF NOT EXISTS `user_online_record_%s`
(
    `id`           BIGINT PRIMARY KEY NOT NULL auto_increment,
    `user_id`      BIGINT             NOT NULL,
    `conn_id`      BIGINT             NOT NULL,
    `platform`     VARCHAR(10)        NOT NULL,
    `online_time`  BIGINT             NOT NULL,
    `offline_time` BIGINT             NOT NULL,
    UNIQUE INDEX `USER_ONLINE_STATUS_U_IDX` (`user_id`, `conn_id`)
);