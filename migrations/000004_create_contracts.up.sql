CREATE TABLE IF NOT EXISTS contracts (
    id UUID PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT'
        CHECK (status IN ('DRAFT', 'PENDING_REVIEW', 'APPROVED', 'ACTIVE', 'EXPIRED', 'TERMINATED')),
    amount_cents BIGINT NOT NULL DEFAULT 0,
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    clauses JSONB NOT NULL DEFAULT '[]',
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    owner_id UUID NOT NULL REFERENCES users(id),
    party_id UUID REFERENCES parties(id),
    template_id UUID REFERENCES templates(id),
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_contracts_owner ON contracts(owner_id);
CREATE INDEX idx_contracts_status ON contracts(status);
CREATE INDEX idx_contracts_party ON contracts(party_id);
CREATE INDEX idx_contracts_created_at ON contracts(created_at DESC);
