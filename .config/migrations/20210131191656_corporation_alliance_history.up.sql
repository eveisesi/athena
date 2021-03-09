CREATE TABLE `corporation_alliance_history` (
    `corporation_id` INT UNSIGNED NOT NULL,
    `record_id` INT UNSIGNED NOT NULL,
    `alliance_id` INT UNSIGNED NULL DEFAULT NULL,
    `is_deleted` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `start_date` TIMESTAMP NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`corporation_id`, `record_id`) USING BTREE,
    INDEX `corporation_alliance_history_alliance_id_idx` (`alliance_id`) USING BTREE,
    INDEX `corporation_alliance_history_corporation_id_start_date_idx` (`corporation_id`, `start_date`),
    CONSTRAINT `corporation_alliance_history_alliance_id_alliances_id_foreign` FOREIGN KEY (`alliance_id`) REFERENCES `athena`.`alliances` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;