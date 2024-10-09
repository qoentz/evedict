CREATE TABLE prediction (
                            id UUID PRIMARY KEY,
                            headline VARCHAR(255) NOT NULL,
                            summary TEXT,
                            image_url VARCHAR(255),
                            timestamp TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_prediction_timestamp ON prediction(timestamp);
