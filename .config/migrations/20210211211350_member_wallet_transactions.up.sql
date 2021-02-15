CREATE TABLE `member_wallet_transactions` (
    `member_id` INT UNSIGNED NOT NULL,
    `transaction_id` BIGINT UNSIGNED NOT NULL,
    `journal_ref_id` BIGINT UNSIGNED NOT NULL,
    `client_id` INT UNSIGNED NOT NULL,
    `client_type` VARCHAR(64) NOT NULL,
    `location_id` BIGINT UNSIGNED NOT NULL,
    `location_type` VARCHAR(64) NOT NULL,
    `type_id` INT UNSIGNED NOT NULL,
    `quantity` INT UNSIGNED NOT NULL,
    `unit_price` FLOAT NOT NULL,
    `is_buy` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `is_personal` TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `date` TIMESTAMP NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    PRIMARY KEY (`member_id`, `transaction_id`) USING BTREE,
    INDEX `member_wallet_transactions_journal_ref_id` (`journal_ref_id`) USING BTREE,
    INDEX `member_wallet_transactions_client_id` (`client_id`) USING BTREE,
    INDEX `member_wallet_transactions_client_type` (`client_type`) USING BTREE,
    INDEX `member_wallet_transactions_location_id` (`location_id`) USING BTREE,
    INDEX `member_wallet_transactions_location_type` (`location_type`) USING BTREE,
    INDEX `member_wallet_transactions_type_id` (`type_id`) USING BTREE,
    INDEX `member_wallet_transactions_date` (`date`) USING BTREE,
    CONSTRAINT `member_wallet_transactions_member_id_foreign` FOREIGN KEY (`member_id`) REFERENCES `athena`.`members` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;