DO $$
BEGIN
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON users;
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON categories;
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON review_patterns;
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON review_boxes;
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON review_items;
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON pattern_steps;
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON scheduled_review_dates;
END;
$$;

DROP FUNCTION IF EXISTS set_updated_at();
