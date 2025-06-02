CREATE TABLE scheduled_review_dates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES review_items(id) ON DELETE CASCADE,
    step_number SMALLINT NOT NULL,
    scheduled_date DATE NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(item_id, step_number)
);

CREATE INDEX idx_scheduled_review_dates_scheduled_review_date ON scheduled_review_dates (scheduled_date);