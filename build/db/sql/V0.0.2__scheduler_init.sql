
-- The table of stream job instance.
CREATE TABLE IF NOT EXISTS `stream_job_instance` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Job ID it belongs to
    `job_id` CHAR(20) NOT NULL,

    -- The release version
    `version` CHAR(16) NOT NULL,

    -- Instance ID
    `id` CHAR(20) NOT NULL,

    -- Instance state, 0 => "pending", 1 => "running", 2 => "suspended", 3 => "successful", 4 => "failed"
    `state` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- Job status, 1 => "Deleted", 2 => "Enabled".
    `status` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,

    -- To save the error message when task execute failed.
    `message` TEXT DEFAULT NULL,

    -- Job instance created time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Job instance updated time
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    INDEX mul_list_record_by_space_id(`space_id`),
    INDEX mul_list_record_by_job_id_version(`job_id`, `version`)

) ENGINE=InnoDB COMMENT='The stream job instance info';

