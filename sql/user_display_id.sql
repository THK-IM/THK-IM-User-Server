CREATE TABLE IF NOT EXISTS `user_display_id_%s`
(
    `display_id` varchar(20) NOT NULL COMMENT '显示id',
    `id`         BIGINT      NOT NULL COMMENT 'id',
    UNIQUE INDEX `User_Id_Display_IDX` (`display_id`),
    UNIQUE INDEX `User_Id_IDX` (`id`)
);