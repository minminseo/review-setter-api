DROP TYPE IF EXISTS back_reviewdate_input;

-- 複合型を再作成
CREATE TYPE back_reviewdate_input AS (
    id uuid,
    category_id uuid,
    box_id uuid,
    initial_scheduled_date date,
    scheduled_date date,
    is_completed boolean
);