CREATE TABLE forecast_market (
                                 forecast_id UUID REFERENCES forecast(id) ON DELETE CASCADE,
                                 market_id BIGINT REFERENCES market(id) ON DELETE CASCADE,
                                 PRIMARY KEY (market_id, forecast_id)
);