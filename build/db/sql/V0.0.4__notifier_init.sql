
CREATE TABLE `notification_post` (
    `notification_post_id` VARCHAR(24) NOT NULL,
    `owner` VARCHAR(65) DEFAULT '' NOT NULL ,
    -- email / sms
    `notify_type` VARCHAR(255) DEFAULT '' NOT NULL,
    `title` VARCHAR(1024) NOT NULL,
    `content` TEXT NOT NULL,
    `short_content` TEXT NOT NULL,
    `status` VARCHAR(50) NOT NULL,
    `create_time` BIGINT(20) UNSIGNED NOT NULL,
    `status_time` BIGINT(20) UNSIGNED NOT NULL,
    `email_address` VARCHAR(1024) NOT NULL,

    PRIMARY KEY (`notification_post_id`)
) ENGINE=InnoDB
