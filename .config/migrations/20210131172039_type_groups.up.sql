CREATE TABLE `type_groups` (
	`id` INT UNSIGNED NOT NULL,
	`name` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_unicode_ci',
	`published` TINYINT NOT NULL,
	`category_id` INT UNSIGNED NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `type_groups_category_id_idx` (`category_id`) USING BTREE,
	CONSTRAINT `type_groups_category_id_categories_id_foreign` FOREIGN KEY (`category_id`) REFERENCES `athena`.`type_categories` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;