CREATE TABLE forecast (
                            id UUID PRIMARY KEY,
                            headline VARCHAR(255) NOT NULL,
                            summary TEXT,
                            image_url VARCHAR,
                            category VARCHAR(255) NOT NULL,
                            timestamp TIMESTAMPTZ NOT NULL,
                            CHECK (category IN ('Politics', 'Economy', 'Technology', 'Culture'))
);

CREATE INDEX idx_forecast_timestamp ON forecast(timestamp);
