-- v2sleep.sub_config definition

CREATE TABLE `sub_config`
(
    `id`        int unsigned NOT NULL AUTO_INCREMENT,
    `name`      varchar(100) NOT NULL,
    `type`      tinyint unsigned NOT NULL,
    `data`      blob,
    `create_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `config_UN` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
