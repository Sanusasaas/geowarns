CREATE TABLE IF NOT EXISTS location_checks (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    checked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_location_checks_user ON location_checks (user_id);
