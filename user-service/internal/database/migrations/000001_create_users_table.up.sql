CREATE SCHEMA IF NOT EXISTS user_service;
SET search_path TO user_service;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY kEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
