-- garante schema
CREATE SCHEMA IF NOT EXISTS project;

-- tabela
CREATE TABLE IF NOT EXISTS project.beer_style
(
    id SERIAL PRIMARY KEY,
    style VARCHAR(100) UNIQUE,
    max_temp INTEGER,
    min_temp INTEGER,
    average_temp NUMERIC(10,2)
);

-- índice de busca por estilo (case-insensitive)
CREATE INDEX IF NOT EXISTS beer_style_style_lower_idx
    ON project.beer_style (lower(style));
