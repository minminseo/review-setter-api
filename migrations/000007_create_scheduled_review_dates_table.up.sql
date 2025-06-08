CREATE TABLE review_dates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    box_id UUID REFERENCES review_boxes(id) ON DELETE SET NULL,
    item_id UUID NOT NULL REFERENCES review_items(id) ON DELETE CASCADE,
    step_number SMALLINT NOT NULL,
    initial_scheduled_date DATE NOT NULL,
    scheduled_date DATE NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(item_id, step_number)
);

CREATE INDEX idx_review_dates_scheduled_date ON review_dates (scheduled_date);