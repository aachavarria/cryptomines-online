-- Migration: Add warehouse fields for auto-production system
-- Module 3: Resource Auto-Production
-- Date: 2026-02-07

-- Add warehouse fields to resources table
-- These fields store accumulated production that hasn't been collected yet

-- Step 1: Add new columns with defaults
ALTER TABLE resources
ADD COLUMN warehouse_metal BIGINT NOT NULL DEFAULT 0 CHECK (warehouse_metal >= 0),
ADD COLUMN warehouse_he3 BIGINT NOT NULL DEFAULT 0 CHECK (warehouse_he3 >= 0),
ADD COLUMN warehouse_gold BIGINT NOT NULL DEFAULT 0 CHECK (warehouse_gold >= 0),
ADD COLUMN last_warehouse_update TIMESTAMPTZ NOT NULL DEFAULT now();

-- Step 2: Initialize warehouse values for existing resources
-- Calculate pending production since last collection and store in warehouse

UPDATE resources
SET
    warehouse_metal = LEAST(
        GREATEST(0,
            CAST(EXTRACT(EPOCH FROM (now() - last_collected_at)) / 3600.0 * metal_per_hour AS BIGINT)
        ),
        storage_capacity - metal
    ),
    warehouse_he3 = LEAST(
        GREATEST(0,
            CAST(EXTRACT(EPOCH FROM (now() - last_collected_at)) / 3600.0 * he3_per_hour AS BIGINT)
        ),
        storage_capacity - he3
    ),
    warehouse_gold = LEAST(
        GREATEST(0,
            CAST(EXTRACT(EPOCH FROM (now() - last_collected_at)) / 3600.0 * gold_per_hour AS BIGINT)
        ),
        storage_capacity - gold
    ),
    last_warehouse_update = now()
WHERE last_collected_at IS NOT NULL;

-- Step 3: For resources that have never been collected, set warehouse to 0 and update timestamp
UPDATE resources
SET
    warehouse_metal = 0,
    warehouse_he3 = 0,
    warehouse_gold = 0,
    last_warehouse_update = now()
WHERE last_collected_at IS NULL;

-- Add index on last_warehouse_update for worker queries
CREATE INDEX idx_resources_last_warehouse_update ON resources (last_warehouse_update);

-- Add comment explaining the warehouse system
COMMENT ON COLUMN resources.warehouse_metal IS 'Accumulated metal production waiting to be collected. Auto-updated by warehouse worker.';
COMMENT ON COLUMN resources.warehouse_he3 IS 'Accumulated He3 production waiting to be collected. Auto-updated by warehouse worker.';
COMMENT ON COLUMN resources.warehouse_gold IS 'Accumulated gold production waiting to be collected. Auto-updated by warehouse worker.';
COMMENT ON COLUMN resources.last_warehouse_update IS 'Last time warehouse values were calculated. Used by auto-production worker.';
