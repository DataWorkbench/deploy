CREATE TABLE `user` (
    `user_id` VARCHAR(50) NOT NULL,
    `user_name` TEXT ,
    `password` VARCHAR(255) DEFAULT '' NOT NULL,
    `lang` VARCHAR(16) DEFAULT 'cn' NOT NULL,
    `email` VARCHAR(255) DEFAULT '' NOT NULL,
    `phone` VARCHAR(50) DEFAULT '' NOT NULL,
    `status` VARCHAR(50) NOT NULL,
    `role` VARCHAR(50) DEFAULT 'user' NOT NULL,
    `currency` VARCHAR(10) DEFAULT 'cny' NOT NULL,
    `gravatar_email` VARCHAR(255) DEFAULT '' NOT NULL,
    `privilege` INTEGER DEFAULT 1 NOT NULL,
    `zones` TEXT ,
    `regions` TEXT ,
    `create_time` BIGINT(20) UNSIGNED NOT NULL,
    `status_time` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`user_id`)
) ENGINE=InnoDB;

CREATE TABLE `access_key` (
    `access_key_id` VARCHAR(50) NOT NULL,
    `access_key_name` TEXT ,
    `secret_access_key` VARCHAR(255) DEFAULT '' NOT NULL,
    `description` TEXT,
    `owner` VARCHAR(50) NOT NULL,
    `status` VARCHAR(50) NOT NULL,
    `ip_white_list` TEXT,
    `create_time` BIGINT(20) UNSIGNED NOT NULL,
    `status_time` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`access_key_id`)
) ENGINE=InnoDB;
