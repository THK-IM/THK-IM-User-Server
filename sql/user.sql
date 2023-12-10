CREATE TABLE IF NOT EXISTS `user_%s`
(
    `id`          BIGINT primary key not null auto_increment,
    `display_id`  BIGINT             NOT NULL COMMENT '显示id',
    `nickname`    varchar(100)       NOT NULL COMMENT '名称',
    `sex`         Tinyint            NOT NULL DEFAULT 0 COMMENT '性别',
    `birthday`    BIGINT             COMMENT '出生日期，时间戳',
    `avatar`      TEXT COMMENT '头像',
    `qr_code`     TEXT COMMENT '二维码',
    `update_time` BIGINT             NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    `create_time` BIGINT             NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    UNIQUE INDEX `User_Display_IDX` (`display_id`)
);