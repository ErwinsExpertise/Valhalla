ALTER TABLE `characters`
  ADD COLUMN `partnerID` int(11) DEFAULT NULL AFTER `vipTeleportRocks`,
  ADD COLUMN `marriageItemID` int(11) DEFAULT NULL AFTER `partnerID`,
  ADD COLUMN `divorceUntil` bigint(20) NOT NULL DEFAULT '0' AFTER `marriageItemID`,
  ADD KEY `partnerID` (`partnerID`),
  ADD CONSTRAINT `characters_ibfk_partner` FOREIGN KEY (`partnerID`) REFERENCES `characters` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

CREATE TABLE IF NOT EXISTS `marriages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `husbandID` int(11) NOT NULL,
  `wifeID` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `husbandID` (`husbandID`),
  KEY `wifeID` (`wifeID`),
  CONSTRAINT `marriages_ibfk_1` FOREIGN KEY (`husbandID`) REFERENCES `characters` (`id`) ON DELETE CASCADE,
  CONSTRAINT `marriages_ibfk_2` FOREIGN KEY (`wifeID`) REFERENCES `characters` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
