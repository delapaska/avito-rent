CREATE TABLE Flat (
    id SERIAL PRIMARY KEY,
    house_id INT REFERENCES House(id) ON DELETE CASCADE,
    price INT NOT NULL,  
    rooms INT NOT NULL,  
    status VARCHAR(25) NOT NULL, 
    moderator_id UUID
);


CREATE INDEX idx_flat_house_number 
ON Flat(house_id);