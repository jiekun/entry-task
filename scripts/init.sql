create database entry_task_v2_db;

CREATE TABLE `user_tab` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'username',
    `nickname` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'nickname',
    `profile_pic` varchar(1024) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'user avatar url',
    `password` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'user password (encrypted)',
    `status` tinyint unsigned NOT NULL COMMENT 'user status 0-enabled 1-disabled',
    `ctime` int unsigned NOT NULL COMMENT 'create timestamp',
    `mtime` int unsigned NOT NULL COMMENT 'modify timestamp',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci