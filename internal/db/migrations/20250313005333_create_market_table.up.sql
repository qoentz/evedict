CREATE TABLE market (
                        id BIGINT PRIMARY KEY,
                        question TEXT NOT NULL,
                        outcomes TEXT NOT NULL,
                        outcome_prices TEXT NOT NULL,
                        volume TEXT NOT NULL,
                        image_url VARCHAR
);