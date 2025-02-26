CREATE TABLE forecast_tag (
                              forecast_id UUID NOT NULL REFERENCES forecast(id) ON DELETE CASCADE,
                              tag_id INT NOT NULL REFERENCES tag(id) ON DELETE CASCADE,
                              PRIMARY KEY (forecast_id, tag_id)
);

CREATE INDEX idx_forecast_tag_forecast_id ON forecast_tag(forecast_id);
CREATE INDEX idx_forecast_tag_tag_id ON forecast_tag(tag_id);
