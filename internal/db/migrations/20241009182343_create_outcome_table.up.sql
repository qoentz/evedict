CREATE TABLE outcome (
                         id UUID PRIMARY KEY,
                         forecast_id UUID NOT NULL REFERENCES forecast(id) ON DELETE CASCADE,
                         content TEXT NOT NULL,
                         confidence_level INT CHECK (confidence_level >= 0 AND confidence_level <= 100) NOT NULL
);

CREATE INDEX idx_outcome_forecast_id ON outcome(forecast_id);