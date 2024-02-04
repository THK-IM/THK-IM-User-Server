CREATE TABLE IF NOT EXISTS `user_%s`
(
    `id`          BIGINT PRIMARY KEY NOT NULL COMMENT 'id',
    `display_id`  varchar(20)        NOT NULL COMMENT '显示id',
    `nickname`    varchar(32) COMMENT '名称',
    `sex`         Tinyint COMMENT '性别',
    `phone`       varchar(20) COMMENT '手机号',
    `birthday`    BIGINT COMMENT '出生日期，时间戳',
    `avatar`      TEXT COMMENT '头像',
    `qrcode`      TEXT COMMENT '二维码',
    `update_time` BIGINT             NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    `create_time` BIGINT             NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    UNIQUE INDEX `User_Display_IDX` (`display_id`)
);