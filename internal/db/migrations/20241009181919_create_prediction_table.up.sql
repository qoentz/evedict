CREATE TABLE forecast (
                            id UUID PRIMARY KEY,
                            headline VARCHAR(255) NOT NULL,
                            summary TEXT,
                            image_url VARCHAR(255),
                            timestamp TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_forecast_timestamp ON forecast(timestamp);
