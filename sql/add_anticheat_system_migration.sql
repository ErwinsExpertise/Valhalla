-- Anti-Cheat and Ban System Tables
-- Migration to add comprehensive anti-cheat and ban management

-- Table for storing active and historical bans
DROP TABLE IF EXISTS `bans`;
CREATE TABLE `bans` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `accountID` INT(10) UNSIGNED DEFAULT NULL,
  `characterID` INT(11) DEFAULT NULL,
  `ipAddress` VARCHAR(45) DEFAULT NULL,
  `hwid` VARCHAR(20) DEFAULT NULL COMMENT 'Hardware ID (machine ID)',
  `banType` ENUM('temporary', 'permanent') NOT NULL DEFAULT 'temporary',
  `banTarget` ENUM('character', 'account', 'ip', 'hwid') NOT NULL DEFAULT 'account',
  `reason` TEXT NOT NULL,
  `issuedBy` VARCHAR(255) DEFAULT NULL,
  `issuedByGM` TINYINT(1) NOT NULL DEFAULT 0,
  `isActive` TINYINT(1) NOT NULL DEFAULT 1,
  `banStartTime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `banEndTime` TIMESTAMP NULL DEFAULT NULL,
  `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_account` (`accountID`, `isActive`),
  KEY `idx_character` (`characterID`, `isActive`),
  KEY `idx_ip` (`ipAddress`, `isActive`),
  KEY `idx_hwid` (`hwid`, `isActive`),
  KEY `idx_active_bans` (`isActive`, `banEndTime`),
  CONSTRAINT `bans_fk_account` FOREIGN KEY (`accountID`) REFERENCES `accounts` (`accountID`) ON DELETE CASCADE,
  CONSTRAINT `bans_fk_character` FOREIGN KEY (`characterID`) REFERENCES `characters` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Table for tracking temporary ban count for escalation
DROP TABLE IF EXISTS `ban_escalation`;
CREATE TABLE `ban_escalation` (
  `accountID` INT(10) UNSIGNED NOT NULL,
  `tempBanCount` INT(11) NOT NULL DEFAULT 0,
  `lastTempBanTime` TIMESTAMP NULL DEFAULT NULL,
  `permanentBanIssued` TINYINT(1) NOT NULL DEFAULT 0,
  `updatedAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`accountID`),
  CONSTRAINT `ban_escalation_fk_account` FOREIGN KEY (`accountID`) REFERENCES `accounts` (`accountID`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Table for logging all violation events
DROP TABLE IF EXISTS `violation_logs`;
CREATE TABLE `violation_logs` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `accountID` INT(10) UNSIGNED NOT NULL,
  `characterID` INT(11) NOT NULL,
  `ipAddress` VARCHAR(45) DEFAULT NULL,
  `violationType` VARCHAR(64) NOT NULL,
  `violationCategory` ENUM('combat', 'movement', 'inventory', 'economy', 'skill', 'packet') NOT NULL,
  `severity` ENUM('low', 'medium', 'high', 'critical') NOT NULL DEFAULT 'medium',
  `detectionDetails` TEXT,
  `mapID` INT(11) DEFAULT NULL,
  `timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `actionTaken` VARCHAR(128) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_account_time` (`accountID`, `timestamp`),
  KEY `idx_character_time` (`characterID`, `timestamp`),
  KEY `idx_violation_type` (`violationType`, `timestamp`),
  KEY `idx_category` (`violationCategory`, `timestamp`),
  CONSTRAINT `violation_logs_fk_account` FOREIGN KEY (`accountID`) REFERENCES `accounts` (`accountID`) ON DELETE CASCADE,
  CONSTRAINT `violation_logs_fk_character` FOREIGN KEY (`characterID`) REFERENCES `characters` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Table for tracking rolling window violation counts
DROP TABLE IF EXISTS `violation_counters`;
CREATE TABLE `violation_counters` (
  `accountID` INT(10) UNSIGNED NOT NULL,
  `characterID` INT(11) NOT NULL,
  `violationType` VARCHAR(64) NOT NULL,
  `count` INT(11) NOT NULL DEFAULT 1,
  `windowStart` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `lastViolation` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`accountID`, `characterID`, `violationType`),
  KEY `idx_window` (`violationType`, `windowStart`),
  KEY `idx_last_violation` (`lastViolation`),
  CONSTRAINT `violation_counters_fk_account` FOREIGN KEY (`accountID`) REFERENCES `accounts` (`accountID`) ON DELETE CASCADE,
  CONSTRAINT `violation_counters_fk_character` FOREIGN KEY (`characterID`) REFERENCES `characters` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Add HWID column to accounts table for tracking hardware IDs
ALTER TABLE `accounts` ADD COLUMN `hwid` VARCHAR(20) DEFAULT NULL COMMENT 'Hardware ID (machine ID)' AFTER `lastIP`;
ALTER TABLE `accounts` ADD INDEX `idx_hwid` (`hwid`);
