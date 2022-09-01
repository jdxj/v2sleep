-- v2sleep.`section` definition

CREATE TABLE `section`
(
    `id`        int unsigned NOT NULL AUTO_INCREMENT,
    `sig`       varchar(500) NOT NULL,
    `create_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
