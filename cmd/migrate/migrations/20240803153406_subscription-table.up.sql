CREATE TABLE Subscriptions (
    id SERIAL PRIMARY KEY,
    house_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL
);
