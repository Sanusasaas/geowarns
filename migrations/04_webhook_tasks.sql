CREATE TABLE IF NOT EXISTS webhook_tasks (
    id SERIAL PRIMARY KEY,
    incident_id INTEGER NOT NULL REFERENCES incidents(id),
    user_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    payload JSONB,
    attempts INTEGER NOT NULL DEFAULT 0,
    next_attempt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_webhook_tasks_status ON webhook_tasks (status);
CREATE INDEX IF NOT EXISTS idx_webhook_tasks_incident ON webhook_tasks (incident_id);
CREATE INDEX IF NOT EXISTS idx_webhook_tasks_user ON webhook_tasks (user_id);
