-- ============================================================================
-- Phase 4.5: PvP Travel & Protection Systems
-- Tables: pending_attacks
-- Columns: players.space_points, players.max_space_points, players.sp_last_reset
-- ============================================================================

-- =========================
-- 1. SPACE POINTS (SP) on players
-- =========================
ALTER TABLE players ADD COLUMN space_points INTEGER NOT NULL DEFAULT 20;
ALTER TABLE players ADD COLUMN max_space_points INTEGER NOT NULL DEFAULT 20;
ALTER TABLE players ADD COLUMN sp_last_reset TIMESTAMPTZ NOT NULL DEFAULT now();

-- =========================
-- 2. PENDING ATTACKS
-- =========================
CREATE TABLE pending_attacks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attacker_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    defender_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    defender_planet_id UUID NOT NULL REFERENCES planets(id) ON DELETE CASCADE,
    fleet_ids TEXT[] NOT NULL,
    status TEXT NOT NULL DEFAULT 'traveling'
        CHECK (status IN ('traveling', 'resolved', 'cancelled')),
    travel_seconds INTEGER NOT NULL,
    depart_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    arrival_at TIMESTAMPTZ NOT NULL,
    return_at TIMESTAMPTZ,
    combat_report_id UUID,
    sp_cost INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_pending_attacks_traveling ON pending_attacks (arrival_at)
    WHERE status = 'traveling';
CREATE INDEX idx_pending_attacks_defender ON pending_attacks (defender_id)
    WHERE status = 'traveling';
CREATE INDEX idx_pending_attacks_attacker ON pending_attacks (attacker_id)
    WHERE status IN ('traveling', 'resolved');
