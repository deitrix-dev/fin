CREATE TABLE `accounts`
(
    `id`   char(36)    NOT NULL,
    `name` varchar(50) NOT NULL,
    PRIMARY KEY (`id`)
) CHARSET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

CREATE TABLE `recurring_payments`
(
    `id`        char(36)     NOT NULL,
    `name`      varchar(127) NOT NULL,
    `enabled`   bool         NOT NULL,
    `debt`      bool         NOT NULL,
    `schedules` json         NOT NULL,
    PRIMARY KEY (`id`)
) CHARSET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

CREATE TABLE `payments`
(
    `id`                   char(36) NOT NULL,
    `date`                 datetime NOT NULL,
    `amount`               int      NOT NULL,
    `account_id`           char(36) NOT NULL,
    `recurring_payment_id` char(36) NULL,
    PRIMARY KEY (`id`),
    INDEX `payments___fk` (`account_id`),
    INDEX `payments_date_index` (`date`),
    INDEX `payments_recurring_payments_id_fk` (`recurring_payment_id`),
    CONSTRAINT `payments___fk` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT `payments_recurring_payments_id_fk` FOREIGN KEY (`recurring_payment_id`) REFERENCES `recurring_payments` (`id`) ON UPDATE CASCADE ON DELETE SET NULL
) CHARSET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;