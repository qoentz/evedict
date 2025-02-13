CREATE TABLE outcome (
                         id UUID PRIMARY KEY,
                         divination_id UUID NOT NULL REFERENCES divination(id) ON DELETE CASCADE,
                         content TEXT NOT NULL,
                         confidence_level INT CHECK (confidence_level >= 0 AND confidence_level <= 100) NOT NULL
);

CREATE INDEX idx_outcome_divination_id ON outcome(divination_id);