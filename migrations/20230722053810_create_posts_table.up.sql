create table if not exists `posts` (
	`id` bigint unsigned not null auto_increment,
	`user_id` bigint unsigned not null,
	`caption` varchar(1000) not null default '',
	`filepath` varchar(1000),
	`created_at` datetime default NOW(),
	primary key (`id`),
	key (`user_id`)
);
