CREATE TABLE IF NOT EXISTS House (
    "id" SERIAL PRIMARY KEY,
    "address" VARCHAR(255) NOT NULL,
    "year" INT NOT NULL,
    "developer" VARCHAR(255),
    "created_at" TIMESTAMP,
    "updated_at" TIMESTAMP
);