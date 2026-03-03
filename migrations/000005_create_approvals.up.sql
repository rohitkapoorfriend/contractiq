CREATE TABLE IF NOT EXISTS approvals (
    id UUID PRIMARY KEY,
    contract_id UUID NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
    reviewer_id UUID NOT NULL REFERENCES users(id),
    decision VARCHAR(20) NOT NULL DEFAULT 'PENDING'
        CHECK (decision IN ('PENDING', 'APPROVED', 'REJECTED')),
    comment TEXT NOT NULL DEFAULT '',
    decided_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_approvals_contract ON approvals(contract_id);
CREATE INDEX idx_approvals_reviewer ON approvals(reviewer_id);
