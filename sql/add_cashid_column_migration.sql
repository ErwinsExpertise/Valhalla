-- Migration to add cashID column to account_cashshop_storage_items table
-- and cashID/cashSN columns to items table
-- This handles existing databases that were created with the old schema

-- Add cashID to cash shop storage items
ALTER TABLE account_cashshop_storage_items 
ADD COLUMN cashID BIGINT(20) DEFAULT NULL AFTER itemID;

-- Add cashID and cashSN to regular inventory items
-- (so items moved from cash shop storage to inventory retain their tracking info)
ALTER TABLE items
ADD COLUMN cashID BIGINT(20) DEFAULT NULL AFTER creatorName,
ADD COLUMN cashSN INT(11) DEFAULT NULL AFTER cashID;

-- Note: Existing items will have NULL cashID/cashSN values
-- Cash shop storage items with NULL cashID will have IDs auto-generated when loaded
-- (see the Load method in cashshop/storage.go which handles NULL cashID values)
-- Regular inventory items with NULL values are not from cash shop, so this is expected
