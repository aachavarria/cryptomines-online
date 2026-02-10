-- Ship Instances table
-- Individual ship instances for granular tracking (recycling, commanders, etc.)
-- Replaces the quantity-based ships table system

CREATE TABLE ship_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    ship_design_id UUID NOT NULL REFERENCES ship_designs(id) ON DELETE CASCADE,
    hull_type_id INTEGER NOT NULL REFERENCES hull_types(id),
    name TEXT, -- Optional custom name for the ship
    current_shield INTEGER NOT NULL DEFAULT 0,
    current_structure INTEGER NOT NULL DEFAULT 0,
    is_damaged BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ship_instances_player ON ship_instances(player_id);
CREATE INDEX idx_ship_instances_design ON ship_instances(ship_design_id);
CREATE INDEX idx_ship_instances_hull ON ship_instances(hull_type_id);

-- Update fleet_stacks to reference ship_instance_id instead of ship_count
ALTER TABLE fleet_stacks
ADD COLUMN ship_instance_id UUID REFERENCES ship_instances(id) ON DELETE CASCADE;

-- Create index for faster fleet stack lookups
CREATE INDEX idx_fleet_stacks_instance ON fleet_stacks(ship_instance_id);

COMMENT ON TABLE ship_instances IS 'Individual ship instances - each row represents one physical ship';
COMMENT ON COLUMN ship_instances.name IS 'Optional custom name for the ship (e.g., "USS Enterprise")';
COMMENT ON COLUMN ship_instances.current_shield IS 'Current shield HP (for damage tracking)';
COMMENT ON COLUMN ship_instances.current_structure IS 'Current structure HP (for damage tracking)';
COMMENT ON COLUMN ship_instances.is_damaged IS 'Quick flag to identify damaged ships needing repair';
