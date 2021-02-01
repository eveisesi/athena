CREATE TABLE `bloodlines` (
	`id` INT UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_unicode_ci',
	`race_id` INT UNSIGNED NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`),
	INDEX `bloodlines_race_id_idx` (`race_id`),
	CONSTRAINT `bloodlines_race_id_races_id_foreign` FOREIGN KEY (`race_id`) REFERENCES `athena`.`races` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;