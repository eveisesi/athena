CREATE TABLE `structures` (
    `id` BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `solar_system_id` INT UNSIGNED NOT NULL,
    `type_id` INT UNSIGNED NOT NULL,
    `owner_id` INT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `structures_solar_system_id_idx` (`solar_system_id`) USING BTREE,
    INDEX `structures_type_id_idx` (`type_id`) USING BTREE,
    INDEX `structures_owner_id_idx` (`owner_id`) USING BTREE,
    CONSTRAINT `structures_solar_system_id_solar_systems_id_foreign` FOREIGN KEY (`solar_system_id`) REFERENCES `athena`.`solar_systems` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT `structures_type_id_types_id_foreign` FOREIGN KEY (`type_id`) REFERENCES `athena`.`types` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT `structures_owner_id_corporations_id_foreign` FOREIGN KEY (`owner_id`) REFERENCES `athena`.`corporations` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;