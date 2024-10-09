CREATE TABLE outcome (
                         id UUID PRIMARY KEY,
                         prediction_id UUID NOT NULL REFERENCES prediction(id) ON DELETE CASCADE,
                         content TEXT NOT NULL,
                         confidence_level INT CHECK (confidence_level >= 0 AND confidence_level <= 100) NOT NULL
);

CREATE INDEX idx_outcome_prediction_id ON outcome(prediction_id);