BEGIN;

CREATE TABLE IF NOT EXISTS todos(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    complete BOOLEAN DEFAULT FALSE,
    expiring_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

COMMIT;