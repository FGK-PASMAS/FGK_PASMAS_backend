CREATE TABLE IF NOT EXISTS "division" (
    id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE IF NOT EXISTS "passenger" (
    id BIGSERIAL PRIMARY KEY,
    last_name TEXT,
    first_name TEXT,
    weight INTEGER NOT NULL,
    division_id INTEGER,

    FOREIGN KEY(division_id) REFERENCES division(id)
);