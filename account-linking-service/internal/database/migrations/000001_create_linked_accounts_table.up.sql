CREATE SCHEMA IF NOT EXISTS account_service;
SET search_path TO account_service;

CREATE TABLE IF NOT EXISTS linked_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    provider_id UUID NOT NULL,
    account_identifier TEXT NOT NULL,
    encrypted_credential TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (user_id, provider_id)
);

