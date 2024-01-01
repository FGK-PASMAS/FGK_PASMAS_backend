WITH division_checks AS (
    SELECT
        (COUNT(*) = 2) AS columns_ok, -- expecting 2 columns
        EXISTS (
            SELECT FROM
                information_schema.columns
            WHERE
                table_name = 'division' AND
                column_name = 'id' AND
                data_type = 'integer'
        ) AND EXISTS (
            SELECT FROM
                information_schema.columns
            WHERE
                table_name = 'division' AND
                column_name = 'name' AND
                data_type = 'text'
        ) AS structure_ok
    FROM
        information_schema.columns
    WHERE
        table_name = 'division'
),
passenger_checks AS (
    SELECT
        (COUNT(*) = 5) AS columns_ok, -- expecting 5 columns
        EXISTS (
            SELECT FROM
                information_schema.columns
            WHERE
                table_name = 'passenger' AND
                column_name = 'id' AND
                data_type = 'bigint'
        ) AND EXISTS (
            SELECT FROM
                information_schema.columns
            WHERE
                table_name = 'passenger' AND
                column_name = 'last_name' AND
                data_type = 'text'
        ) AND EXISTS (
            SELECT FROM
                information_schema.columns
            WHERE
                table_name = 'passenger' AND
                column_name = 'first_name' AND
                data_type = 'text'
        ) AND EXISTS (
            SELECT FROM
                information_schema.columns
            WHERE
                table_name = 'passenger' AND
                column_name = 'weight' AND
                data_type = 'integer'
        ) AND EXISTS (
            SELECT FROM
                information_schema.columns
            WHERE
                table_name = 'passenger' AND
                column_name = 'division_id' AND
                data_type = 'integer'
        ) AS structure_ok
    FROM
        information_schema.columns
    WHERE
        table_name = 'passenger'
),
fk_checks AS (
    SELECT
        EXISTS (
            SELECT FROM
                information_schema.table_constraints AS tc
                JOIN information_schema.key_column_usage AS kcu
                    ON tc.constraint_name = kcu.constraint_name
                JOIN information_schema.constraint_column_usage AS ccu
                    ON ccu.constraint_name = tc.constraint_name
            WHERE
                tc.constraint_type = 'FOREIGN KEY' AND
                tc.table_name = 'passenger' AND
                kcu.column_name = 'division_id' AND
                ccu.table_name = 'division' AND
                ccu.column_name = 'id'
        ) AS fk_ok
)
SELECT
    (division_checks.columns_ok AND division_checks.structure_ok) AND
    (passenger_checks.columns_ok AND passenger_checks.structure_ok) AND
    fk_checks.fk_ok
FROM
    division_checks, passenger_checks, fk_checks;
