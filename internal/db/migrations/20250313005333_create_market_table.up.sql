CREATE TABLE market (
                        id UUID PRIMARY KEY,
                        forecast_id UUID NOT NULL UNIQUE REFERENCES forecast(id) ON DELETE CASCADE,
                        question TEXT NOT NULL,
                        outcomes TEXT NOT NULL,       -- raw JSON string e.g. "[\"Yes\",\"No\"]"
                        outcome_prices TEXT NOT NULL, -- raw JSON string e.g. "[\"0.115\",\"0.885\"]"
                        volume TEXT NOT NULL,         -- keep as TEXT if you store the raw string from the API
                        image_url VARCHAR,
                        external_id VARCHAR
);