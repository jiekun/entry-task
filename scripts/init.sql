create database entry_task_v2_db;

create table `user_tab` (
    `id` bigint unsigned not null auto_increment,
    `name` varchar(64) collate utf8mb4_unicode_ci not null comment 'username',
    `nickname` varchar(64) collate utf8mb4_unicode_ci not null comment 'nickname',
    `profile_pic` varchar(1024) collate utf8mb4_unicode_ci not null comment 'user avatar url',
    `password` varchar(32) character set utf8mb4 collate utf8mb4_unicode_ci not null comment 'user password (encrypted)',
    `status` tinyint unsigned not null comment 'user status 0-enabled 1-disabled',
    `ctime` int unsigned not null comment 'create timestamp',
    `mtime` int unsigned not null comment 'modify timestamp',
    primary key (`id`),
    unique key `uniq_name` (`name`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci;
