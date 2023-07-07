CREATE TABLE `ppm_cmm_industry` (
  `id` bigint(20) NOT NULL,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' ,
  `cname` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' ,
  `pid` bigint(20) NOT NULL DEFAULT '0' ,
  `is_show` tinyint(4) NOT NULL DEFAULT '2' ,
  `is_default` tinyint(4) NOT NULL DEFAULT '2' ,
  `is_delete` tinyint(4) NOT NULL DEFAULT 2,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `ppm_cmm_industry` VALUES (1, 'agriculture', '农业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (2, 'forestry', '林业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (3, 'Animal husbandry\n', '牧业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (4, 'Fishery', '渔业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (5, 'mining industry', '采矿业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (6, 'manufacturing', '制造业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (7, 'electric power', '电力', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (8, 'heat', '热力', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (9, 'Gas', '燃气', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (10, 'Water production', '水生产', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (11, 'Supply industry', '供应业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (12, 'Construction industry', '建筑业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (13, 'Transportation industry', '交通运输业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (14, 'Warehousing industry', '仓储业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (15, 'Postal industry', '邮政业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (16, 'Information transfer', '信息传输业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (17, 'software', '软件', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (18, 'Information technology service industry', '信息技术服务业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (19, 'tourism', '旅游', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (20, 'Fúwù yè\n3/5000\nService industry', '服务业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (21, 'manufacturing', '生产制造', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (23, 'Insurance', '保险业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (24, 'Wholesale and retail trade', '批发和零售业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (25, 'Accommodation and catering industry', '住宿和餐饮业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (26, 'Financial industry', '金融业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (27, 'Real estate industry', '房地产业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (28, 'Leasing and business services', '租赁和商务服务业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (29, 'Scientific research and technical services', '科学研究和技术服务业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (30, 'Water, environmental and public facilities management', '水利、环境和公共设施管理业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (31, 'Resident services, repairs and other services', '居民服务、修理和其他服务业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (32, 'education', '教育', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (33, 'Health and social work', '卫生和社会工作', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (34, 'Culture, sports and entertainment', '文化、体育和娱乐业', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (35, 'Public administration, social security and social organization', '公共管理、社会保障和社会组织', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (36, 'International organizations', '国际组织', 0, 1, 2, 2);
INSERT INTO `ppm_cmm_industry` VALUES (37, 'energy', '能源', 0, 1, 2, 2);

ALTER TABLE `ppm_org_department`
ADD COLUMN `path`  varchar(255) NOT NULL DEFAULT "0" COMMENT '父级路径拼接（“,”隔开）' AFTER `parent_id`;

CREATE TABLE `ppm_org_user_invite` (
  `id` bigint(20) NOT NULL,
  `org_id` bigint NOT NULL DEFAULT '0',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '姓名',
  `email` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '邮箱',
  `invite_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '邀请人id',
  `is_register` int NOT NULL DEFAULT '2' COMMENT '是否注册（1已注册2未注册）',
  `last_invite_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '上次邀请时间',
  `creator` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int NOT NULL DEFAULT '1',
  `is_delete` tinyint NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_org_user_invite_org_id` (`org_id`),
  KEY `index_ppm_org_user_invite_create_time` (`create_time`),
  KEY `index_ppm_org_user_invite_invite_user_id` (`invite_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `ppm_orc_config` ADD COLUMN `db_id` bigint NOT NULL DEFAULT "0" COMMENT '数据库id' AFTER `time_difference`;
ALTER TABLE `ppm_orc_config` ADD COLUMN `ds_id` bigint NOT NULL DEFAULT "0" COMMENT '数据源id' AFTER `time_difference`;
ALTER TABLE `ppm_orc_config` ADD COLUMN `dc_id` bigint NOT NULL DEFAULT "0" COMMENT '数据中心id' AFTER `time_difference`;

ALTER TABLE `ppm_org_user_organization` ADD COLUMN `invite_id` bigint NOT NULL DEFAULT "0" COMMENT '邀请人id' AFTER `user_id`;
