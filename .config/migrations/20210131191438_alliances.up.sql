CREATE TABLE `alliances` (
    `id` INT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `ticker` VARCHAR(255) NOT NULL,
    `creator_id` INT UNSIGNED NOT NULL,
    `creator_corporation_id` INT UNSIGNED NOT NULL,
    `executor_corporation_id` INT UNSIGNED NOT NULL,
    `date_founded` TIMESTAMP NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `alliances_creator_id` (`creator_id`) USIGN BTREE,
    INDEX `alliances_creator_corporation_id` (`creator_corporation_id`) USIGN BTREE,
    INDEX `alliances_executor_corporation_id` (`executor_corporation_id`) USIGN BTREE,
    -- CONSTRAINT `alliances_creator_id_characters_id_foreign` FOREIGN KEY (`creator_id`) REFERENCES `athena`.`characters` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT,
    -- CONSTRAINT `alliances_creator_corporation_id_corporations_id_foreign` FOREIGN KEY (`creator_corporation_id`) REFERENCES `athena`.`corporations` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT,
    -- CONSTRAINT `alliances_executor_corporation_id_corporations_id_foreign` FOREIGN KEY (`executor_corporation_id`) REFERENCES `athena`.`alliances` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;