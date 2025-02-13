CREATE TABLE source (
                        id UUID PRIMARY KEY,
                        divination_id UUID NOT NULL REFERENCES divination(id) ON DELETE CASCADE,
                        name VARCHAR(255) NOT NULL,
                        title VARCHAR(255) NOT NULL,
                        url VARCHAR(255) NOT NULL
);

CREATE INDEX idx_source_divination_id ON source(divination_id);