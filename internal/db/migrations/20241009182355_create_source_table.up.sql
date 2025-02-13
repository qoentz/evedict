CREATE TABLE source (
                        id UUID PRIMARY KEY,
                        forecast_id UUID NOT NULL REFERENCES forecast(id) ON DELETE CASCADE,
                        name VARCHAR(255) NOT NULL,
                        title VARCHAR(255) NOT NULL,
                        url VARCHAR(255) NOT NULL
);

CREATE INDEX idx_source_forecast_id ON source(forecast_id);