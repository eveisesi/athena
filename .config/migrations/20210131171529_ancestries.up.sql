CREATE TABLE `ancestries` (
	`id` INT(10) UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`bloodline_id` INT(10) UNSIGNED NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `ancestries_bloodline_id_idx` (`bloodline_id`) USING BTREE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;