CREATE TABLE `mail_headers` (
	`id` INT UNSIGNED NOT NULL,
	`sender_id` INT UNSIGNED NOT NULL,
	`sender_type` ENUM('character', 'corporation', 'mailing_list') NOT NULL,
	`subject` VARCHAR(255) NULL DEFAULT NULL,
	`body` TEXT NULL DEFAULT NULL,
	`is_on_mailing_list` TINYINT(1) NULL DEFAULT '0',
	`mailing_list_id` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
	`is_ready` TINYINT(1) NOT NULL DEFAULT '0',
	`sent` TIMESTAMP NULL DEFAULT NULL,
	`created_at` TIMESTAMP NULL DEFAULT NULL,
	`updated_at` TIMESTAMP NULL DEFAULT NULL,
	PRIMARY KEY (`id`),
	INDEX `mail_headers_sender_id_index` (`sender_id`)
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;