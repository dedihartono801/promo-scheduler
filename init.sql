CREATE DATABASE IF NOT EXISTS sayakaya;

use sayakaya;

CREATE TABLE IF NOT EXISTS `promo_type` (
  `id` int unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` varchar(100) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
    INDEX(`id`)
);

INSERT IGNORE INTO promo_type (name)
VALUES ("Birthday Promo");

CREATE TABLE IF NOT EXISTS `promo` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `promo_type_id` int unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text,
  `code` varchar(100),
  `discount_type` ENUM('fix', 'percentage'),
  `discount` int,
  `maximum_discount_usage` int,
  `minimum_transaction` int,
  `image` varchar(255),
  `user_type` ENUM('all', 'valid_user'),
  `start_date` datetime NOT NULL,
  `end_date` datetime NOT NULL,
  `is_active` ENUM('1', '0'),
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
    INDEX(`id`),
    INDEX(`promo_type_id`),
    FOREIGN KEY (`promo_type_id`) REFERENCES promo_type(`id`)
);

CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `email` varchar(100) NOT NULL,
  `phone` varchar(15) NOT NULL,
  `date_birth` date NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
    INDEX(`id`),
    INDEX(`email`),
    INDEX(`phone`)
);

INSERT IGNORE INTO user (name,email,phone,date_birth)
VALUES ("Agung","agung@gmail.com","08982510066","1996-11-09");

CREATE TABLE IF NOT EXISTS `user_promo` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `promo_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `is_active` ENUM('1', '0'),
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
    INDEX(`id`),
    INDEX(`promo_id`),
    INDEX(`user_id`),
    INDEX(`is_active`),
    FOREIGN KEY (`promo_id`) REFERENCES promo(`id`),
    FOREIGN KEY (`user_id`) REFERENCES user(`id`)
);
