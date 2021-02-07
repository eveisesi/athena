CREATE TABLE `mailing_lists` (
	`id` INT UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`source_character_id` INT UNSIGNED NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`)
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;