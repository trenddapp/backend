-- Wordle Status
CREATE TYPE wordle_status AS ENUM (
    'OPEN',
    'COMPLETE',
    'CANCELED'
);

-- Wordles
CREATE TABLE IF NOT EXISTS wordles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now() ON UPDATE now(),
    user_id UUID NOT NULL,
    status wordle_status NOT NULL DEFAULT 'OPEN',
    solution STRING NOT NULL,
    guesses STRING[]
)
