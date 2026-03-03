CREATE INDEX idx_contracts_title_search ON contracts USING gin(to_tsvector('english', title));
CREATE INDEX idx_contracts_description_search ON contracts USING gin(to_tsvector('english', description));
