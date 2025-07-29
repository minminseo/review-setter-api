DROP TYPE IF EXISTS back_reviewdate_input;

-- 古い定義(000012)で複合型を再作成 (initial_scheduled_dateなし)
CREATE TYPE back_reviewdate_input AS (
    id uuid,
    category_id uuid,
    box_id uuid,
    scheduled_date date,
    is_completed boolean
);