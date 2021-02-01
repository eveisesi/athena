CREATE TABLE `solar_systems` (
    `id` INT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_unicode_ci',
    `constealltion_id` INT UNSIGNED NOT NULL,
    `security_class` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_unicode_ci',
    `security_status` FLOAT NOT NULL,
    `star_id` INT UNSIGNED NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `solar_systems_constellation_id_idx` (`constealltion_id`) USING BTREE,
    INDEX `solar_systems_star_id_idx` (`star_id`) USING BTREE,
    INDEX `solar_systems_security_status` (`security_status`) USING BTREE,
    CONSTRAINT `solar_systems_constellation_id_constellations_id_foreign` FOREIGN KEY (`constealltion_id`) REFERENCES `athena`.`constellations` ON UPDATE RESTRICT ON DELETE RESTRICT
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;