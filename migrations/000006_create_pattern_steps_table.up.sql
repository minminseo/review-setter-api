CREATE TABLE pattern_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pattern_id UUID NOT NULL REFERENCES review_patterns(id) ON DELETE CASCADE,
    step_number SMALLINT NOT NULL,
    interval_days SMALLINT NOT NULL CHECK (interval_days > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(pattern_id, step_number)
);