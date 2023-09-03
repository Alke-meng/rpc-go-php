USE test;

CREATE TABLE IF NOT EXISTS `tbl_crm_9` (
  `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(90) DEFAULT NULL,
  `sex` INT(10) UNSIGNED DEFAULT NULL,
  `number1` VARCHAR(90) DEFAULT NULL,
  `number2` VARCHAR(90) DEFAULT NULL,
  `remark` TEXT DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `number1` (`number1`),
  KEY `name` (`name`)
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `tbl_crm_recycle_9` (
  `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(90) DEFAULT NULL,
  `sex` INT(10) UNSIGNED DEFAULT NULL,
  `number1` VARCHAR(90) DEFAULT NULL,
  `number2` VARCHAR(90) DEFAULT NULL,
  `remark` TEXT DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `number1` (`number1`),
  KEY `name` (`name`)
) ENGINE=INNODB;


CREATE TABLE IF NOT EXISTS `tbl_callee_pool_9` (
  `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `customer_id` INT(10) UNSIGNED NOT NULL,
  `number` VARCHAR(90) NOT NULL,
  `status` TINYINT(4) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  KEY `customer_id` (`customer_id`),
  KEY `number` (`number`,`status`)
) ENGINE=INNODB;