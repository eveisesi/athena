CREATE TABLE `corporation_alliance_history` (
    `id` INT UNSIGNED NOT NULL,
    `record_id` INT UNSIGNED NOT NULL,
    `alliance_id` INT UNSIGNED NULL DEFAULT NULL,
    `is_deleted` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `start_date` TIMESTAMP NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`id`, `record_id`) USIGN BTREE,
    INDEX `corporation_alliance_history_alliance_id_idx` (`alliance_id`) USIGN BTREE,
    CONSTRAINT `corporation_alliance_history_alliance_id_alliances_id_foreign` FOREIGN KEY (`alliance_id`) REFERENCES `athena`.`alliances` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
)