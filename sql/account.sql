CREATE TABLE IF NOT EXISTS `account_%s`
(
    `user_id`     BIGINT      NOT NULL COMMENT '用户id',
    `channel`     varchar(8)  NOT NULL COMMENT '渠道',
    `account`     varchar(64) NOT NULL COMMENT '账号/凭证',
    `password`    varchar(32) COMMENT '密码',
    `update_time` BIGINT      NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    `create_time` BIGINT      NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    INDEX `Account_User_IDX` (`user_id`),
    INDEX `Account_Channel_IDX` (`account`, `channel`)
);