-- Create Users Table
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    gender VARCHAR(50) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    birth_date DATE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('manager', 'employee'))
);

-- -- Create Tickets Table
-- CREATE TABLE tickets (
--     ticket_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     user_id UUID NOT NULL REFERENCES users(user_id),
--     title VARCHAR(255) NOT NULL,
--     content TEXT NOT NULL,
--     priority VARCHAR(50) NOT NULL CHECK (priority IN ('low', 'medium', 'high')),
--     status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected'))
-- );
