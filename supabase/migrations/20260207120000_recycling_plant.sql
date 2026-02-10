-- Recycling Plant system
-- Allows players to scrap ships for resources

CREATE TABLE recycling_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    ship_design_id UUID NOT NULL REFERENCES ship_designs(id) ON DELETE CASCADE,
    metal_gained BIGINT NOT NULL DEFAULT 0,
    he3_gained BIGINT NOT NULL DEFAULT 0,
    gold_gained BIGINT NOT NULL DEFAULT 0,
    duration_seconds INT NOT NULL DEFAULT 60,
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    collected BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_recycling_jobs_player ON recycling_jobs(player_id);
CREATE INDEX idx_recycling_jobs_completion ON recycling_jobs(completed_at) WHERE collected = false;

COMMENT ON TABLE recycling_jobs IS 'Tracks ship recycling jobs in the Recycling Plant';
