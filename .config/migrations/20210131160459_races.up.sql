CREATE TABLE `races` (
	`id` INT UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_unicode_ci',
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`)
) COLLATE='utf8mb4_unicode_ci' ENGINE=InnoDB;
