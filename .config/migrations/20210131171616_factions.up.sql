CREATE TABLE `factions` (
	`id` INT(10) UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_unicode_ci',
	`race_id` INT(10) UNSIGNED NULL DEFAULT NULL,
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `factions_race_id_idx` (`race_id`) USING BTREE,
	CONSTRAINT `factions_race_id_races_id_foreign` FOREIGN KEY (`race_id`) REFERENCES `athena`.`races` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;