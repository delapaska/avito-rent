CREATE TABLE Users (
    user_id UUID PRIMARY KEY,       
    email VARCHAR(255) UNIQUE NOT NULL,   
    password VARCHAR(255) NOT NULL,    
    user_type VARCHAR(50) NOT NULL        
);

 
CREATE INDEX idx_users_email ON Users(email);