create table if not exists `comments`(
	`id` bigint unsigned not null auto_increment,
	`post_id` bigint unsigned not null,
	`user_id` bigint unsigned not null,
	`text` varchar(1000) not null default '',
	`created_at` datetime default NOW(),
	primary key (`id`),
	key (`post_id`, `created_at` desc)
);
