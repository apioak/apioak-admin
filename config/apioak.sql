SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for oak_certificates
-- ----------------------------
DROP TABLE IF EXISTS `oak_certificates`;
CREATE TABLE `oak_certificates` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Certificate id',
  `sni` varchar(150) NOT NULL DEFAULT '' COMMENT 'SNI',
  `certificate` text NOT NULL DEFAULT '' COMMENT 'Certificate content',
  `private_key` text NOT NULL DEFAULT '' COMMENT 'Private key content',
  `enable` tinyint(1) unsigned NOT NULL DEFAULT 2 COMMENT 'Certificate enable  1:on  2:off',
  `expired_at` timestamp NULL DEFAULT NULL COMMENT 'Expiration time',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  KEY `IDX_SNI` (`sni`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Certificates';

-- ----------------------------
-- Table structure for oak_plugin_configs
-- ----------------------------
DROP TABLE IF EXISTS `oak_plugin_configs`;
CREATE TABLE `oak_plugin_configs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Plugin config id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT 'Plugin config name',
  `type` tinyint(2) NOT NULL DEFAULT 0 COMMENT 'Plugin relation type 1:service  2:router',
  `target_id` char(20) NOT NULL DEFAULT '' COMMENT 'Target id',
  `plugin_res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Plugin res id',
  `plugin_key` varchar(20) NOT NULL DEFAULT '' COMMENT 'Plugin key',
  `config` text NOT NULL COMMENT 'Plugin configuration',
  `enable` tinyint(1) unsigned NOT NULL DEFAULT 2 COMMENT 'Plugin config enable  1:on  2:off',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Plugin Configs';

-- ----------------------------
-- Table structure for oak_plugins
-- ----------------------------
DROP TABLE IF EXISTS `oak_plugins`;
CREATE TABLE `oak_plugins` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Plugin id',
  `plugin_key` varchar(20) NOT NULL DEFAULT '' COMMENT 'Plugin key',
  `icon` varchar(50) NOT NULL DEFAULT '' COMMENT 'Plugin icon',
  `type` tinyint(2) NOT NULL DEFAULT 0 COMMENT 'Plugin type',
  `description` varchar(200) NOT NULL DEFAULT '' COMMENT 'Plugin description',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  UNIQUE KEY `UNIQ_KEY` (`plugin_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Plugins';

-- ----------------------------
-- Table structure for oak_routers
-- ----------------------------
DROP TABLE IF EXISTS `oak_routers`;
CREATE TABLE `oak_routers` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Router id',
  `service_res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Service id',
  `upstream_res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Upstream id',
  `router_name` varchar(50) NOT NULL DEFAULT '' COMMENT 'Router name',
  `request_methods` varchar(150) NOT NULL DEFAULT '' COMMENT 'Request method',
  `router_path` varchar(200) NOT NULL DEFAULT '' COMMENT 'Routing path',
  `enable` tinyint(1) unsigned NOT NULL DEFAULT 2 COMMENT 'Router enable  1:on  2:off',
  `release` tinyint(1) unsigned NOT NULL DEFAULT 1 COMMENT 'Service release status 1:unpublished  2:to be published  3:published',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Routers';

-- ----------------------------
-- Table structure for oak_service_domains
-- ----------------------------
DROP TABLE IF EXISTS `oak_service_domains`;
CREATE TABLE `oak_service_domains` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Domain id',
  `service_res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Service id',
  `domain` varchar(50) NOT NULL DEFAULT '' COMMENT 'Domain name',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  UNIQUE KEY `UNIQ_SERVICE_ID_DOMAIN` (`service_res_id`,`domain`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Service domains';

-- ----------------------------
-- Table structure for oak_services
-- ----------------------------
DROP TABLE IF EXISTS `oak_services`;
CREATE TABLE `oak_services` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Service id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT 'Service name',
  `protocol` tinyint(1) unsigned NOT NULL DEFAULT 1 COMMENT 'Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS',
  `enable` tinyint(1) unsigned NOT NULL DEFAULT 2 COMMENT 'Service enable  1:on  2:off',
  `release` tinyint(1) unsigned NOT NULL DEFAULT 1 COMMENT 'Service release status 1:unpublished  2:to be published  3:published',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  KEY `IDX_NAME` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Services';

-- ----------------------------
-- Table structure for oak_upstream_nodes
-- ----------------------------
DROP TABLE IF EXISTS `oak_upstream_nodes`;
CREATE TABLE `oak_upstream_nodes` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Service node id',
  `upstream_res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Upstream id',
  `node_ip` varchar(60) NOT NULL DEFAULT '' COMMENT 'Node IP',
  `ip_type` tinyint(1) unsigned NOT NULL DEFAULT 1 COMMENT 'IP Type  1:IPV4  2:IPV6',
  `node_port` smallint(6) unsigned NOT NULL DEFAULT 0 COMMENT 'Node port',
  `node_weight` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT 'Node weight',
  `health` tinyint(1) unsigned NOT NULL DEFAULT 1 COMMENT 'Health type  1:HEALTH  2:UNHEALTH',
  `health_check` tinyint(1) unsigned NOT NULL DEFAULT 2 COMMENT 'Health check  1:on  2:off',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  UNIQUE KEY `UNIQ_UPSTREAM_ID_NODE_IP_PORT` (`upstream_res_id`,`node_ip`,`node_port`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Upstream nodes';

-- ----------------------------
-- Table structure for oak_upstreams
-- ----------------------------
DROP TABLE IF EXISTS `oak_upstreams`;
CREATE TABLE `oak_upstreams` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Upstream id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT 'Upstream name',
  `algorithm` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT 'Load balancing algorithm  1:round robin  2:chash',
  `connect_timeout` int(10) unsigned NOT NULL DEFAULT 1 COMMENT 'Connect timeout',
  `write_timeout` int(10) unsigned NOT NULL DEFAULT 1 COMMENT 'Write timeout',
  `read_timeout` int(10) unsigned NOT NULL DEFAULT 1 COMMENT 'Read timeout',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Upstreams';

-- ----------------------------
-- Table structure for oak_user_tokens
-- ----------------------------
DROP TABLE IF EXISTS `oak_user_tokens`;
CREATE TABLE `oak_user_tokens` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'User tokenID',
  `token` text NOT NULL DEFAULT '' COMMENT 'Token',
  `user_email` varchar(80) NOT NULL DEFAULT '' COMMENT 'Email',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Update time',
  `expired_at` timestamp NULL DEFAULT NULL COMMENT 'Expired time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  UNIQUE KEY `UNIQ_USER_EMAIL` (`user_email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User token';

-- ----------------------------
-- Table structure for oak_users
-- ----------------------------
DROP TABLE IF EXISTS `oak_users`;
CREATE TABLE `oak_users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'User iD',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT 'User name',
  `password` char(32) NOT NULL DEFAULT '' COMMENT 'Password',
  `email` varchar(80) NOT NULL DEFAULT '' COMMENT 'Email',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  UNIQUE KEY `UNIQ_EMAIL` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Users';

SET FOREIGN_KEY_CHECKS = 1;
