package database

var structureStatements = `
    CREATE TABLE IF NOT EXISTS "division" (
        id SERIAL PRIMARY KEY,
        name TEXT,
        passenger_capacity INTEGER NOT NULL
    );

    CREATE TABLE IF NOT EXISTS "passenger" (
        id BIGSERIAL PRIMARY KEY,
        last_name TEXT,
        first_name TEXT,
        weight INTEGER NOT NULL
    );
`

var seedDivisionStatements = `
    INSERT INTO division(name, passenger_capacity) VALUES ('Segelflug', 1);
    INSERT INTO division(name, passenger_capacity) VALUES ('Motorsegler', 1);
    INSERT INTO division(name, passenger_capacity) VALUES ('Motorflug', 3);
`
