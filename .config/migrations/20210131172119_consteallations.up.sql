CREATE TABLE `constellations` (
    `id` INT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_unicode_ci',
    `region_id` INT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `constellations_region_id_idx` (`region_id`) USING BTREE,
    CONSTRAINT `constellations_region_id_regions_id_foreign` FOREIGN KEY (`region_id`) REFERENCES `athena`.`regions` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;