CREATE TABLE source (
                        id UUID PRIMARY KEY,
                        prediction_id UUID NOT NULL REFERENCES prediction(id) ON DELETE CASCADE,
                        name VARCHAR(255) NOT NULL,
                        title VARCHAR(255) NOT NULL,
                        url VARCHAR(255) NOT NULL
);

CREATE INDEX idx_source_prediction_id ON source(prediction_id);