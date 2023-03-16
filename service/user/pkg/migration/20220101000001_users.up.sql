-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now() ON UPDATE now(),
    address STRING NOT NULL UNIQUE,
    balance INT NOT NULL
);

-- Nonces
CREATE TABLE IF NOT EXISTS nonces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    value STRING NOT NULL
);
