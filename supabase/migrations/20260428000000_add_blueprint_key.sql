-- =============================================================================
-- Add blueprint_key column to blueprints table.
--
-- Why: quest reward JSON references blueprints by a stable lowercase key
-- (e.g. "estrella", "ship_reinforcement_facility"). The previous lookup tried
-- to match against hull_types.name, which has tier suffixes (estrella_i),
-- producing "Failed to find blueprint estrella: sql: no rows in result set".
--
-- Derivation rules:
--   - Hull blueprints  -> hull_types.name with trailing "_i" stripped
--                         (estrella_i -> estrella, air_wanderer_i -> air_wanderer)
--   - Module blueprints -> module_types.name (already a stable lowercase key)
-- =============================================================================

ALTER TABLE blueprints ADD COLUMN blueprint_key TEXT;

UPDATE blueprints b
SET blueprint_key = regexp_replace(h.name, '_i$', '')
FROM hull_types h
WHERE b.hull_type_id = h.id;

UPDATE blueprints b
SET blueprint_key = m.name
FROM module_types m
WHERE b.module_type_id = m.id;

ALTER TABLE blueprints ALTER COLUMN blueprint_key SET NOT NULL;
ALTER TABLE blueprints ADD CONSTRAINT blueprints_blueprint_key_unique UNIQUE (blueprint_key);
CREATE INDEX idx_blueprints_blueprint_key ON blueprints (blueprint_key);
