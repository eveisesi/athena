ALTER TABLE
ADD
    CONSTRAINT `corporations_home_station_id_stations_id_foreign` FOREIGN KEY (`home_station_id`) REFERENCES `athena`.`stations` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT,
ADD
    CONSTRAINT `corporations_faction_id_factions_id_foreign` FOREIGN KEY (`faction_id`) REFERENCES `athena`.`factions` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT;