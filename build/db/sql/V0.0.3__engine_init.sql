
CREATE TABLE IF NOT EXISTS `nad_ipam` (
    -- flink_id.
    `flink_id` VARCHAR(32) NOT NULL,

    -- vxnet_id.
    `vxnet_id` VARCHAR(32) NOT NULL,

    -- The ipam subnet from vxnet.router.ip_network.
    `subnet` VARCHAR(22) NOT NULL,

    -- The ipam gateway from vxnet.router.manager_ip.
    `gateway` VARCHAR(18) NOT NULL,

    -- The start ip for flink cluster.
    `start_ip` VARCHAR(18) NOT NULL,

    -- The end ip for flink cluster.
    `end_ip` VARCHAR(18) NOT NULL,

    -- The status of flink cluster deleted: 0, active: 1.
    `status` TINYINT(1) NOT NULL DEFAULT 1,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`flink_id`)
) ENGINE=InnoDB COMMENT='The ipam info in NetworkAttachmentDefinition.';
