-- World Chat system
-- Allows players to communicate globally

CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    channel VARCHAR(50) NOT NULL DEFAULT 'world',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_chat_messages_channel_time ON chat_messages(channel, created_at DESC);
CREATE INDEX idx_chat_messages_player ON chat_messages(player_id);

COMMENT ON TABLE chat_messages IS 'Stores chat messages for world and alliance channels';

-- Rate limiting table
CREATE TABLE chat_rate_limits (
    player_id UUID PRIMARY KEY REFERENCES players(id) ON DELETE CASCADE,
    last_message_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    message_count INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON TABLE chat_rate_limits IS 'Tracks player message rate for anti-spam';
