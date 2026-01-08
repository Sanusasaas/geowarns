CREATE TABLE IF NOT EXISTS incident_stats (
    id SERIAL PRIMARY KEY,
    incident_id INTEGER NOT NULL REFERENCES incidents(id),
    total_checks INTEGER NOT NULL DEFAULT 0,
    total_near INTEGER NOT NULL DEFAULT 0,
    unique_users INTEGER NOT NULL DEFAULT 0,
    last_check_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (incident_id)
);

CREATE INDEX IF NOT EXISTS idx_incident_stats_incident ON incident_stats (incident_id);
