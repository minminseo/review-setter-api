CREATE TYPE target_weight_enum AS ENUM ('heavy', 'normal', 'light', 'unset');

CREATE TABLE review_patterns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    target_weight target_weight_enum NOT NULL DEFAULT 'unset',
    registered_at TIMESTAMPTZ NOT NULL,
    edited_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);