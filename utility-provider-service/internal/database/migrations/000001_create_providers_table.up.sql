CREATE TABLE providers (
    id UUID PRIMARY kEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    api_url VARCHAR(255) NOT NULL,
    authentication_type VARCHAR(50) NOT NULL,
    api_key VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
