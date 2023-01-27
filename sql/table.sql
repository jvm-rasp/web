create database jrasp;
use jrasp;
DROP TABLE IF EXISTS `rbac_permission`;
CREATE TABLE `rbac_permission`  (
                                    `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                                    `parent_id` int(10) NULL DEFAULT 0,
                                    `name` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
                                    `icon` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
                                    `url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
                                    `hide_in_menu` tinyint(1) UNSIGNED NULL DEFAULT 1 COMMENT '是否在导航菜单中',
                                    `hide_children_in_menu` tinyint(1) UNSIGNED NULL DEFAULT 1,
                                    `priority` int(10) UNSIGNED NULL DEFAULT 0,
                                    `created_at` datetime NULL DEFAULT NULL,
                                    `updated_at` datetime NULL DEFAULT NULL,
                                    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC;

DROP TABLE IF EXISTS `rbac_role`;
CREATE TABLE `rbac_role`  (
                              `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                              `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
                              `code` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
                              `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
                              PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC;

DROP TABLE IF EXISTS `rbac_role_permission`;
CREATE TABLE `rbac_role_permission`  (
                                         `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                                         `role_id` int(10) UNSIGNED NULL DEFAULT 0,
                                         `permission_id` int(10) NULL DEFAULT NULL,
                                         PRIMARY KEY (`id`) USING BTREE,
                                         UNIQUE INDEX `role_id,permission_id`(`role_id`, `permission_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 482 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC;

DROP TABLE IF EXISTS `rbac_user`;
CREATE TABLE `rbac_user`  (
                              `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                              `avatar` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '',
                              `username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
                              `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
                              `mobile` char(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
                              `status` tinyint(1) UNSIGNED NULL DEFAULT 1,
                              `created_at` datetime NULL DEFAULT NULL,
                              `updated_at` datetime NULL DEFAULT NULL,
                              PRIMARY KEY (`id`) USING BTREE,
                              UNIQUE INDEX `username_idx`(`username`) USING BTREE,
                              UNIQUE INDEX `mobile_idx`(`mobile`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC;

DROP TABLE IF EXISTS `rbac_user_action_log`;
CREATE TABLE `rbac_user_action_log`  (
                                         ` id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                                         `user_id` int(10) UNSIGNED NULL DEFAULT NULL COMMENT '用户id',
                                         `operation` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '操作内容',
                                         `item_id` int(10) UNSIGNED NULL DEFAULT 0 COMMENT '相关id',
                                         `ip` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
                                         `created_at` datetime NULL DEFAULT NULL,
                                         PRIMARY KEY (` id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC;

DROP TABLE IF EXISTS `rbac_user_role`;
CREATE TABLE `rbac_user_role`  (
                                   `id` int(10) NOT NULL AUTO_INCREMENT,
                                   `user_id` int(10) UNSIGNED NULL DEFAULT 0,
                                   `role_id` int(10) NULL DEFAULT NULL,
                                   PRIMARY KEY (`id`) USING BTREE,
                                   UNIQUE INDEX `user_id,role_id`(`user_id`, `role_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC;

DROP TABLE IF EXISTS `rasp_host`;
CREATE TABLE `rasp_host`
(
    `id`               bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '递增id',
    `host_name`        varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实例名称',
    `ip`               varchar(64) COLLATE utf8mb4_unicode_ci NULL COMMENT 'ip地址',
    `agent_mode`       varchar(16) COLLATE utf8mb4_unicode_ci  NOT NULL DEFAULT 'disable' COMMENT 'agent接入模式:disable;static;dynamic',
    `agent_info`       text COLLATE utf8mb4_unicode_ci                  DEFAULT NULL COMMENT 'java agent信息',
    `config_id`        int unsigned NULL COMMENT '策略id',
    `heartbeat_time`   timestamp NULL DEFAULT NULL COMMENT '最近一次的心跳时间',
    `install_dir`      varchar(255) COLLATE utf8mb4_unicode_ci NULL COMMENT '可执行文件安装的绝对路径',
    `version`          varchar(32) COLLATE utf8mb4_unicode_ci NULL COMMENT 'rasp 版本',
    `exe_file_hash`    varchar(255) COLLATE utf8mb4_unicode_ci NULL COMMENT '可执行文件的hash',
    `os_type`          varchar(16) COLLATE utf8mb4_unicode_ci           DEFAULT NULL COMMENT '实例操作系统类型:darwin、windowns、linux等',

    `total_mem`        int DEFAULT 0 COMMENT '实例总内存,单位GB',
    `cpu_counts`       int DEFAULT 0 COMMENT '实例逻辑cpu数量',
    `free_disk`        int DEFAULT 0 COMMENT '实例可用磁盘容量,单位GB',

    `build_date_time`  varchar(64) COLLATE utf8mb4_unicode_ci          DEFAULT NULL COMMENT '代码信息：编译时间',
    `build_git_branch` varchar(64) COLLATE utf8mb4_unicode_ci          DEFAULT NULL COMMENT '代码信息：分支',
    `build_git_commit` varchar(64) COLLATE utf8mb4_unicode_ci          DEFAULT NULL COMMENT '代码信息：commit',
    `nacos_info`       varchar(256) COLLATE utf8mb4_unicode_ci          DEFAULT NULL COMMENT 'nacos初始化信息，json格式',
    `agent_config_update_time` timestamp NULL DEFAULT NULL COMMENT 'agent 配置更新时间',

    `status`           int                                              DEFAULT 1 COMMENT '是否禁用:0,永久下线；1,在线',
    `create_time`      timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`      timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_host_name` (`host_name`),
    KEY                `idx_heartbeat_time` (`heartbeat_time`)
) ENGINE = InnoDB AUTO_INCREMENT = 1000 CHARACTER
SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `rasp_java_info`;
CREATE TABLE `rasp_java_info`
(
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '递增id',
    `host_name`       varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实例名称',
    `cmdline_info`    text COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '进程命令行信息，文本格式',
    `start_time`      timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '启动时间',
    `pid`             int unsigned NOT NULL COMMENT '进程号',
    `status`          varchar(64) COLLATE utf8mb4_unicode_ci NULL COMMENT '注入状态',
    `dependency_info` text COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '依赖信息，json字符串格式',
    `create_time`     timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`     timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_host_name_and_pid` (`host_name`,`pid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `rasp_config`;
CREATE TABLE `rasp_config`
(
    `id`             int unsigned NOT NULL AUTO_INCREMENT COMMENT '递增id',
    `config_name`    varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '策略名称',
    `message`    	 varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '策略描述',
    `config_content` text COLLATE utf8mb4_unicode_ci         NOT NULL COMMENT '策略内容',
    `tag`            varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '策略标签',
    `username`       varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '创建的用户名',
    `status`         int unsigned DEFAULT 1 COMMENT '是否禁用:0,禁用；1,启用',
    `create_time`    timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`    timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_config_name` (`config_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `rasp_module`;
CREATE TABLE `rasp_module`
(
    `id`                     int unsigned NOT NULL AUTO_INCREMENT COMMENT '递增id',
    `name`            varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '插件名称',
    `url`             varchar(1024) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '插件下载地址',
    `hash`            varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '插件hash',
    `type`                   varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模块类型：http、rce、file_read、sql',
    `middleware_name`    	 varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '中间件名称：tomcat、jeety',
    `middleware_version`     varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '中间件版本:8.x、9.x',
    `version`         varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '插件版本:1.0.4、1.0.5',
    `tag`                    varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '策略标签',
    `parameters`     longtext COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '插件参数模版：json字符串',
    `username`               varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '创建的用户名',
    `message`               varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '创建的用户名',
    `status`                 int unsigned DEFAULT 1 COMMENT '是否禁用:0,禁用；1,启用',
    `create_time`    timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`    timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


DROP TABLE IF EXISTS `attack_info`;
CREATE TABLE `attack_info`
(
    `id`                 bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '递增id',
    `host_name`          varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实例名称',
    `local_ip`           varchar(64) COLLATE utf8mb4_unicode_ci NULL COMMENT '被攻击者ip地址',
    `remote_ip`          varchar(64) COLLATE utf8mb4_unicode_ci NULL COMMENT '攻击者ip地址',
    `attack_type`        varchar(64) COLLATE utf8mb4_unicode_ci NULL COMMENT '攻击类型',
    `check_type`         varchar(256) COLLATE utf8mb4_unicode_ci NULL COMMENT '检测算法类型以及信息',
    `is_blocked`         tinyint(1) COLLATE utf8mb4_unicode_ci NULL COMMENT '是否阻断：0,放行；1：阻断',
    `level`              tinyint COLLATE utf8mb4_unicode_ci               DEFAULT 0 COMMENT '危险等级：0～100整数',
    `handle_status`      tinyint COLLATE utf8mb4_unicode_ci      NOT NULL DEFAULT 1 COMMENT '处理状态，未处理1(初始化状态)、处理中2(点击查看详情)、已经处理3',
    `handle_result`      text COLLATE utf8mb4_unicode_ci                  DEFAULT NULL COMMENT '处理结果：确认漏洞、误拦截、忽略',
    `stack_trace`        text COLLATE utf8mb4_unicode_ci                  DEFAULT NULL COMMENT '调用栈：json字符串格式',
    `http_method`        varchar(16) COLLATE utf8mb4_unicode_ci NULL COMMENT 'http请求类型：get、post、put',
    `request_protocol`   varchar(16) COLLATE utf8mb4_unicode_ci NULL COMMENT 'http请求协议：rpc、http,dubbo',
    `request_uri`        varchar(256) COLLATE utf8mb4_unicode_ci NULL COMMENT '请求路径',
    `request_parameters` varchar(2048) COLLATE utf8mb4_unicode_ci NULL COMMENT '请求参数：json字符串格式',
    `attack_parameters`  varchar(2048) COLLATE utf8mb4_unicode_ci NULL COMMENT '攻击参数：json字符串格式',
    `cookies`            varchar(2048) COLLATE utf8mb4_unicode_ci NULL COMMENT '请求cookie：json字符串格式',
    `header`             varchar(2048) COLLATE utf8mb4_unicode_ci NULL COMMENT '请求header：json字符串格式',
    `body`               text COLLATE utf8mb4_unicode_ci  DEFAULT NULL COMMENT 'http body 信息',
    `tag`                varchar(128) COLLATE utf8mb4_unicode_ci NULL COMMENT '危险等级的文字描述:严重、高危、中危、低危',
    `attack_time`        timestamp NULL DEFAULT NULL COMMENT '攻击时间',
    `create_time`        timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`        timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY                  `idx_host_name` (`host_name`),
    KEY                  `idx_local_ip` (`local_ip`),
    KEY                  `idx_remote_ip` (`remote_ip`),
    KEY                  `idx_handle_status` (`handle_status`),
    KEY                  `idx_attack_time` (`attack_time`),
    KEY                  `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `service_info`;
CREATE TABLE `service_info`
(
    `id`                 int unsigned NOT NULL AUTO_INCREMENT COMMENT '递增id',
    `service_name`       varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '服务名称',
    `description`        varchar(1024) COLLATE utf8mb4_unicode_ci NULL COMMENT '服务描述',
    `language`           varchar(64) COLLATE utf8mb4_unicode_ci NULL COMMENT '开发语言：java、go、c++、python等',
    `owners`             varchar(1024) COLLATE utf8mb4_unicode_ci NULL COMMENT '服务负责人',
    `node_number`        int unsigned DEFAULT 0 COMMENT '节点数量，默认为0',
    `is_public`          tinyint(1) COLLATE utf8mb4_unicode_ci DEFAULT 0 COMMENT '是否公开：0,公开，所有可见、可操作；1：私有，仅服务负责人可以看到，默认公开',
    `is_inner`           tinyint(1) COLLATE utf8mb4_unicode_ci DEFAULT 0 COMMENT '是否外网可以访问：0,外网；1：内网，默认外网',
    `level`              tinyint(1) COLLATE utf8mb4_unicode_ci DEFAULT 0 COMMENT '服务等级：0，非核心服务；1，核心服务，默认非核心服务',
    `low_period`         int COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 4 COMMENT '低峰期：0～23小时，默认为4:00～4:59',
    `organization`        varchar(128) COLLATE utf8mb4_unicode_ci NULL COMMENT '服务的组织：信息安全、交易中心等',
    `tag`                varchar(128) COLLATE utf8mb4_unicode_ci NULL COMMENT '标签，自定义标签',
    `create_time`        timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`        timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE  KEY          `unique_idx_service_name` (`service_name`),
    KEY                  `idx_create_time` (`create_time`),
    KEY                  `idx_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2022-05-12 (同步客户1) 原始表结构已经修改
-- ALTER TABLE attack_info ADD COLUMN body text COLLATE utf8mb4_unicode_ci  DEFAULT NULL COMMENT 'http body 信息' after header

-- 2022-07-31 (同步客户2) 原始表结构已经修改
-- ALTER TABLE rasp_host ADD COLUMN build_date_time varchar(64) COLLATE utf8mb4_unicode_ci  DEFAULT NULL COMMENT '代码信息：编译时间' after free_disk;
-- ALTER TABLE rasp_host ADD COLUMN build_git_branch varchar(64) COLLATE utf8mb4_unicode_ci  DEFAULT NULL COMMENT '代码信息：分支' after build_date_time;
-- ALTER TABLE rasp_host ADD COLUMN build_git_commit varchar(64) COLLATE utf8mb4_unicode_ci  DEFAULT NULL COMMENT '代码信息：commit' after build_git_branch;

-- 2022-08-18 (同步客户1) 原始表结构已经修改
-- ALTER TABLE rasp_host ADD COLUMN nacos_info varchar(256) COLLATE utf8mb4_unicode_ci  DEFAULT NULL COMMENT 'nacos初始化信息，json格式' after build_git_commit;
-- ALTER TABLE rasp_host ADD COLUMN agent_config_update_time timestamp NULL DEFAULT NULL COMMENT 'agent 配置更新时间' after nacos_info;
