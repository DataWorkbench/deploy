-- Workspace Info
CREATE TABLE IF NOT EXISTS `workspace` (
    -- Workspace ID, unique within a region
    `id` CHAR(20) NOT NULL,

     -- User ID of workspace owner
    `owner` VARCHAR(65) NOT NULL,

    -- Workspace name, unique within a region
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- Workspace description
    `desc` VARCHAR(1024) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- Workspace status, 1 => "enabled", 2 => "disabled", 3 => "deleted"
    `status` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE INDEX unique_space_name (`owner`, `name`)

) ENGINE=InnoDB COMMENT='The workspace info';

-- Workspace Operation Audit
CREATE TABLE IF NOT EXISTS `audit` (
    -- Only used to query sort by
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,

    -- The user id of who execute this operation.
    `user_id`  VARCHAR(65) NOT NULL,

    -- The workspace id
    `space_id` VARCHAR(24),

    -- The type of operation permission,  1 => "Write", 2 => "Read".
    `perm_type` TINYINT(1) UNSIGNED NOT NULL,

    -- The operation of user behavior.
    `api_name` VARCHAR(128) NOT NULL,

    -- The operation state, 1 => "Success", 2 => "Failed".
    `state` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- Timestamp of time of when accessed.
    `created` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    -- Index to list all records of the specified user.
    INDEX mul_list_audit_by_user_id(`user_id`),
    -- Index to lists all records of the specified workspace id.
    INDEX mul_list_audit_by_space_id(`space_id`)
) ENGINE=InnoDB COMMENT='The workspace operation opaudit record';

-- Workspace Member
CREATE TABLE IF NOT EXISTS `member` (
    -- Only used to query sort by
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,

    -- Workspace ID that the member belongs to.
    `space_id` VARCHAR(24) NOT NULL,

    -- Use account user-id as member id.
    `user_id` VARCHAR(65) NOT NULL,

    -- Member status, 1 => "Normal" 2 => "Deleted".
    `status` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,

    -- Workspace description
    `desc` VARCHAR(1024) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- The id lists of system role. Multiple id separated by commas, eg: "ros-1,ros-2".
    `system_role_ids` VARCHAR(256) NOT NULL,

    -- The id lists of custom role. Multiple id separated by commas, eg: "roc-1,roc-2"
    -- A member can have up to 100 custom roles.
    `custom_role_ids` VARCHAR(2048) NOT NULL,

    -- User ID of created this member.
    `created_by` VARCHAR(65) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE key unique_user_id(`space_id`, `user_id`)

) ENGINE =Innodb COMMENT ='The workspace member';

-- The table of user quota.
CREATE TABLE IF NOT EXISTS `user_quota` (
    -- The user id.
    `user_id` VARCHAR(65) NOT NULL,

    `workspace` JSON,
    `member` JSON,
    `custom_role` JSON,
    `stream_job` JSON,
    `sync_job` JSON,
    `data_source` JSON,
    `udf` JSON,
    `file` JSON,
    `flink_cluster` JSON,
    `network` JSON,
    PRIMARY KEY (`user_id`)
) ENGINE=InnoDB COMMENT='The quota limit for user level.';

-- The table of workspace quota.
CREATE TABLE IF NOT EXISTS `workspace_quota` (
    `space_id` CHAR(20) NOT NULL,

    `workspace` JSON,
    `member` JSON,
    `custom_role` JSON,
    `stream_job` JSON,
    `sync_job` JSON,
    `data_source` JSON,
    `udf` JSON,
    `file` JSON,
    `flink_cluster` JSON,
    `network` JSON,
    PRIMARY KEY (`space_id`)
) ENGINE=InnoDB COMMENT='The quota limit for workspace level.';

-- The table of user quota.
    CREATE TABLE IF NOT EXISTS `member_quota` (
    `space_id` CHAR(20) NOT NULL,
    -- The user id of member.
    `user_id` VARCHAR(65) NOT NULL,

    `workspace` JSON,
    `member` JSON,
    `custom_role` JSON,
    `stream_job` JSON,
    `sync_job` JSON,
    `data_source` JSON,
    `udf` JSON,
    `file` JSON,
    `flink_cluster` JSON,
    `network` JSON,
    PRIMARY KEY (`space_id`, `user_id`)
) ENGINE=InnoDB COMMENT='The quota limit for member level.';


-- The table of stream job.
CREATE TABLE IF NOT EXISTS `stream_job` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Job ID, unique within a region
    `id` CHAR(20) NOT NULL,

    -- The workflow version id
    `version` CHAR(16) NOT NULL,

    -- PID is the parent id(directory). pid is "" means root(`/`)
    `pid` CHAR(20) NOT NULL,

    -- IsDirectory represents this job whether a directory.
    `is_directory` BOOL,

    -- Job Name, Unique within a workspace.
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- Job description
    `desc` VARCHAR(1024) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- Job type, 0 = "NoType", 1 => "StreamOperator" 2 => "StreamSQL" 3 => "StreamJAR" 4 => "StreamPython"
    `type` TINYINT(1) UNSIGNED NOT NULL,

    -- Workspace status, 1 => "deleted", 2 => "enabled"
    `status` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,

    -- User ID of created this job.
    `created_by` VARCHAR(65) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`, `version`),
    UNIQUE KEY unique_job_name (`space_id`, `version`, `name`)

) ENGINE=InnoDB COMMENT='The stream job info.';

-- The table of stream job property.
CREATE TABLE IF NOT EXISTS `stream_job_property` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Job ID it belongs to
    `id` CHAR(20) NOT NULL,

    -- Release version, unique
    `version` CHAR(16) NOT NULL,

    -- The job code that format with JSON.
    `code` JSON,

    -- The environment parameters that format with JSON.
    `args` JSON,

    -- The schedule property that format with JSON.
    `schedule` JSON,

    PRIMARY KEY (`id`, `version`)

) ENGINE=InnoDB COMMENT='The meta of stream workflow.';

-- The table of stream job release.
CREATE TABLE IF NOT EXISTS `stream_job_release` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Job ID it belongs to
    `id` CHAR(20) NOT NULL,

    -- The release version
    `version` CHAR(16) NOT NULL,

    -- Job Name, Unique within a workspace
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- Job type, 1 => "StreamOperator" 2 => "StreamSQL" 3 => "StreamJAR" 4 => "StreamPython" 5 => "StreamScala"
    `type` TINYINT(1) UNSIGNED NOT NULL,

    -- Release status, 1 => "Active", 2 => "Suspended", 3 => "Deleted",
    `status` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,

    -- Job release description
    `desc` VARCHAR(1024) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- User ID of release this job.
    `created_by` VARCHAR(65) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    INDEX mul_list_record_by_space_id(`space_id`)

) ENGINE=InnoDB COMMENT='The release latest info';

-- The table of stream job versions.
-- create table table_name_new like table_name_old;
CREATE TABLE IF NOT EXISTS `stream_job_version` like `stream_job`;

-- The table of stream job meta version.
CREATE TABLE IF NOT EXISTS `stream_job_property_version` like `stream_job_property`;

# -- The table of monitor rule.
# CREATE TABLE IF NOT EXISTS `monitor_rule` (
#     -- Workspace id it belongs to
#     `space_id` CHAR(20) NOT NULL,
#
#     -- The rule id.
#     `id` CHAR(20) NOT NULL,
#
#     -- The rule name.
#     -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
#     -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
#     `name` VARCHAR(149) NOT NULL,
#
#     -- Monitor status, 1 => "deleted", 2 => "enabled", 3 => "disabled"
#     `status` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,
#
#     -- The object unit, 1 => "workspace" 2 => "job'
#     `unit` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,
#
#     -- The object text.
#     `text` VARCHAR(1024) NOT NULL,
#
#     -- The trigger conditions. 3 => "retrying" 6 => "timeout"  7 => "succeed", 8 => "failed",
#     `trigger` TINYINT(1) UNSIGNED NOT NULL,
#
#     -- The alarm times. 1 ~ 99
#     `alarm_times` TINYINT(1) UNSIGNED NOT NULL,
#
#     -- The alarm interval. 1 ~ 30
#     `alarm_interval` TINYINT(1) UNSIGNED NOT NULL,
#
#     -- The alarm type. "sms, email"
#     `alarm_type` VARCHAR(32) NOT NULL,
#
#     -- The free time. "00:01,03:00".
#     `free_time` VARCHAR(16) NOT NULL,
#
#     -- The alarm receiver. "usr-111111,usr-22222".
#     `receiver` VARCHAR(256) NOT NULL,
#
#     PRIMARY KEY (`id`),
#     UNIQUE KEY (`space_id`, `name`)
#
# ) ENGINE=InnoDB COMMENT='The monitor rule';

create table if not exists data_source (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- DataSource ID
    `id` CHAR(20) NOT NULL,

    -- unique in a workspace.
    -- The max length of use set is 64. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 85 (64 + 20 + 1)
    `name` VARCHAR(85) NOT NULL,

    -- DataSource description.
    `desc` varchar(256),

    -- Type, 1->MySQL 2->PostgreSQL 3->Kafka 4->S3 5->ClickHouse 6->Hbase 7->Ftp 8->HDFS
    `type` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- URL of data source settings..
    `url` JSON,

    -- Status, 1 => "Delete", 2 => "enabled", 3 => "disabled"
    `status` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- User ID of created this data source.
    `created_by` VARCHAR(65) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE KEY unique_source_name (`space_id`, `name`)
) ENGINE=InnoDB COMMENT='The data source info';

create table if not exists data_source_connection (
    -- Only used to query sort by
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,

    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- DataSource ID
    `source_id` CHAR(20) NOT NULL,

    -- Network ID
    `network_id` CHAR(20) NOT NULL,

    -- Status, 1 => "Delete", 2 => "Enabled"
    `status` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- result, 1-=> success 2 => failed
    `result` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- Message is the reason when connection failure.
    `message` VARCHAR(1024) NOT NULL,

    -- Use time. unit in ms.
    `elapse` INT,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    INDEX mul_query_with_space_id(`space_id`),
    INDEX mul_query_with_source_id(`source_id`)

) ENGINE=InnoDB COMMENT='The connection info for datasource';

CREATE TABLE IF NOT EXISTS `network` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Network ID, unique within a region
    `id` CHAR(20) NOT NULL,

    -- Network Name, Unique within a workspace
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- VPC's route_id.
    `router_id` VARCHAR(32) NOT NULL,

    -- VPC's vxnet_id.
    `vxnet_id` VARCHAR(32) NOT NULL,

    -- The user-id of created this network.
    `created_by` VARCHAR(128) NOT NULL,

    -- The cluster status. 1 => "deleted" 2 => "Enabled"
    `status` TINYINT(1) UNSIGNED NOT NULL DEFAULT 1,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE KEY unique_network_name (`space_id`, `name`)

) ENGINE=InnoDB COMMENT='The network info.';

CREATE TABLE IF NOT EXISTS `flink_cluster` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Cluster ID, unique within a region
    `id` CHAR(20) NOT NULL,

    -- Cluster Name, Unique within a workspace
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- The flink version.
    `version` VARCHAR(63) NOT NULL,

    -- The cluster status. 1 => "deleted" 2 => "running" 3 => "stopped" 4 => "starting" 5 => "exception" 6 => "Arrears"
    `status` TINYINT(1) UNSIGNED NOT NULL,

    -- Flink task number for TaskManager. Is required, Min 1, Max ?
    `task_num` INT,

    -- Flink JobManager's cpu and memory. 1CU = 1C + 4GB. Is required, Min 0.5, Max 8
    `job_cu` FLOAT,

    -- Flink TaskManager's cpu and memory. 1CU = 1C + 4GB. Is required, Min 0.5, Max 8
    `task_cu` FLOAT,

    -- Network config.
    `network_id` CHAR(20) NOT NULL,

    -- Config of host aliases
    `host_aliases` JSON,

    -- Flink config.
    `config` JSON,

    -- The user-id of created this flink cluster.
    `created_by` VARCHAR(128) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE KEY unique_flink_cluster_name(`space_id`, `name`)

) ENGINE=InnoDB COMMENT='The flink cluster info.';

CREATE TABLE IF NOT EXISTS `udf` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- UDF ID, unique within a region
    `id` CHAR(20) NOT NULL,

    -- UDF Name, unique within a space.
    -- The max length of use set is 64. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 85 (64 + 20 + 1)
    `name` VARCHAR(85) NOT NULL,

    -- Description, describe this UDF.
    `desc` VARCHAR(256) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- UDF status, Optional Values: 1 => "deleted", 2 => "enabled"
    `status` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- UDF Type; Optional Values: 1=>UDF, 2=>UDTF 3=> UDTTF.
    `type` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- UDF language; Optional Values: 1 => Scala 2=> Java 3=> Python
    `language` TINYINT(1) UNSIGNED DEFAULT 0 NOT NULL,

    -- The id of resource. Used with language of JAVA.
    `file_id` CHAR(20) NOT NULL,

    -- The code. Used with language of Python and Scala.
    `code` text NOT NULL,

    -- usage example for this udf
    `usage_sample` varchar(2048),

    -- Who created this udf.
    `created_by` varchar(65),

    -- Timestamp of create time.
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time.
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    UNIQUE KEY (`space_id`, `name`)
) ENGINE=InnoDB COMMENT='The udf schema';

CREATE TABLE IF NOT EXISTS `file`
(
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- File ID, unique within a region
    `id` CHAR(20) NOT NULL,

    -- Id of Parent Directory. pid is "" means root(`/`).
    `pid` CHAR(20) NOT NULL,

    -- IsDirectory represents this job whether a directory.
    `is_directory` BOOL,

    -- File Name, Unique within a workspace
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- File description.
    `desc` VARCHAR(1024) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- File status. 1 => "deleted", 2 => "enabled"
    `status` TINYINT(1) UNSIGNED NOT NULL,

    -- File Size.
    `size` BIGINT(50),

    -- MD5 value of file data encoded in hexadecimal.
    `etag` CHAR(32),

    -- The version of this file.
    `version` CHAR(16) NOT NULL,

    -- Who created this file.
    `created_by` varchar(65),

    -- Timestamp of create time.
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time.
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    INDEX unique_file_name (`space_id`, `pid`, `name`)
) ENGINE=InnoDB COMMENT='The file schema';

-- Table for for describes dependencies between modules.
CREATE TABLE IF NOT EXISTS `binding` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- module_id represents which resources are bound to this module.
    `module_id` CHAR(20) NOT NULL,

    -- module_version is the version of module.
    -- This filed maybe empty.
    `module_version` CHAR(16) NOT NULL,

    -- resource_id represents the module bound resources.
    `resource_id` CHAR(20) NOT NULL,

    -- resource_version is the version of resource.
    -- Notice: Reserved field, unused on present.
    `resource_version` CHAR(16) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`space_id`, `module_id`, `module_version`, `resource_id`, `resource_version`)
#     INDEX mul_index_with_module_id(`module_id`, `module_version`),
#     INDEX mul_index_with_resource_id(`resource_id`, `resource_version`)

) ENGINE=InnoDB COMMENT='describes dependencies between modules.';



-- The table of sync job.
CREATE TABLE IF NOT EXISTS `sync_job` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Job ID, unique within a region
    `id` CHAR(20) NOT NULL,

    -- The workflow version id
    `version` CHAR(16) NOT NULL,

    -- PID is the parent id(directory). pid is "" means root(`/`)
    `pid` CHAR(20) NOT NULL,

    -- IsDirectory represents this job whether a directory.
    `is_directory` BOOL,

    -- Job Name, Unique within a workspace.
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- Job description
    `desc` VARCHAR(1024) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- Job type, 0 = "NoType", 1 => "StreamOperator" 2 => "StreamSQL" 3 => "StreamJAR" 4 => "StreamPython" 5 => "StreamScala"
    `type` TINYINT(1) UNSIGNED NOT NULL,

    -- Workspace status, 1 => "deleted", 2 => "enabled"
    `status` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,

    -- User ID of created this job.
    `created_by` VARCHAR(65) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`, `version`),
    UNIQUE KEY unique_job_name (`space_id`, `version`, `name`)

) ENGINE=InnoDB COMMENT='The sync job info.';

-- The table of sync job property.
CREATE TABLE IF NOT EXISTS `sync_job_property` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Job ID it belongs to
    `id` CHAR(20) NOT NULL,

    -- Release version, unique
    `version` CHAR(16) NOT NULL,

    -- The environment conf that format with JSON.
    `conf` JSON,

    -- The schedule property that format with JSON.
    `schedule` JSON,

    PRIMARY KEY (`id`, `version`)

) ENGINE=InnoDB COMMENT='The meta of sync workflow.';

-- The table of sync job release.
CREATE TABLE IF NOT EXISTS `sync_job_release` (
    -- Workspace id it belongs to
    `space_id` CHAR(20) NOT NULL,

    -- Job ID it belongs to
    `id` CHAR(20) NOT NULL,

    -- The release version
    `version` CHAR(16) NOT NULL,

    -- Job Name, Unique within a workspace
    -- The max length of use set is 128. The system will be auto rename to <name>.<id> when deleted.
    -- Thus the VARCHAR should be define as 149 (128 + 20 + 1)
    `name` VARCHAR(149) NOT NULL,

    -- Job type, 1 => "StreamOperator" 2 => "StreamSQL" 3 => "StreamJAR" 4 => "StreamPython" 5 => "StreamScala"
    `type` TINYINT(1) UNSIGNED NOT NULL,

    -- Release status, 1 => "Active", 2 => "Suspended", 3 => "Deleted",
    `status` TINYINT(1) UNSIGNED DEFAULT 1 NOT NULL,

    -- Job release description
    `desc` VARCHAR(1024) CHARACTER SET utf8mb4 DEFAULT '' NOT NULL,

    -- User ID of release this job.
    `created_by` VARCHAR(65) NOT NULL,

    -- Timestamp of create time
    `created` BIGINT(20) UNSIGNED NOT NULL,

    -- Timestamp of update time, Update when some changed, default value should be same as "created"
    `updated` BIGINT(20) UNSIGNED NOT NULL,

    PRIMARY KEY (`id`),
    INDEX mul_list_record_by_space_id(`space_id`)

) ENGINE=InnoDB COMMENT='The sync job release latest info';

-- The table of sync job versions.
-- create table table_name_new like table_name_old;
CREATE TABLE IF NOT EXISTS `sync_job_version` like `sync_job`;

-- The table of sync job meta version.
CREATE TABLE IF NOT EXISTS `sync_job_property_version` like `sync_job_property`;