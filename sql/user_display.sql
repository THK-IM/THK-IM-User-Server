CREATE TABLE IF NOT EXISTS `user_display_%s`
(
    `display_id` varchar(20) PRIMARY KEY NOT NULL COMMENT '显示id',
    `id`         BIGINT                  NOT NULL COMMENT 'id'
);