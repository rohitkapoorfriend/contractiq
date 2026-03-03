CREATE TABLE IF NOT EXISTS parties (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    email VARCHAR(255) NOT NULL,
    party_type VARCHAR(20) NOT NULL CHECK (party_type IN ('ORGANIZATION', 'INDIVIDUAL')),
    company VARCHAR(200) NOT NULL DEFAULT '',
    phone VARCHAR(20) NOT NULL DEFAULT '',
    address VARCHAR(500) NOT NULL DEFAULT '',
    created_by UUID NOT NULL REFERENCES users(id),
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_parties_created_by ON parties(created_by);
