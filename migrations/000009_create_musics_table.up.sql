CREATE TYPE reviewdate_input AS (
    id uuid,
    category_id uuid,
    box_id uuid,
    scheduled_date date,
    is_completed boolean
);
