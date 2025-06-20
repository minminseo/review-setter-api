-- name: CreateItem :exec
INSERT INTO
    review_items (
        id,
        user_id,
        category_id,
        box_id,
        pattern_id,
        name,
        detail,
        learned_date,
        is_Finished,
        registered_at,
        edited_at
    )
VALUES (
    sqlc.arg(id),
    sqlc.arg(user_id),
    sqlc.arg(category_id),
    sqlc.arg(box_id),
    sqlc.arg(pattern_id),
    sqlc.arg(name),
    sqlc.arg(detail),
    sqlc.arg(learned_date),
    sqlc.arg(is_Finished),
    sqlc.arg(registered_at),
    sqlc.arg(edited_at)
    );

-- 新規一括挿入時と、一括更新時に使う
-- name: CreateReviewDates :copyfrom
INSERT INTO
    review_dates (
        id,
        user_id,
        category_id,
        box_id,
        item_id,
        step_number,
        initial_scheduled_date,
        scheduled_date,
        is_completed
    ) VALUES (
        sqlc.arg(id),
        sqlc.arg(user_id),
        sqlc.arg(category_id),
        sqlc.arg(box_id),
        sqlc.arg(item_id),
        sqlc.arg(step_number),
        sqlc.arg(initial_scheduled_date),
        sqlc.arg(scheduled_date),
        sqlc.arg(is_completed)
    );


-- 学習日変更など、どういうリクエストなのかを判定するために使う
-- name: GetItemByID :one
SELECT
    id,
    user_id,
    category_id,
    box_id,
    pattern_id,
    name,
    detail,
    learned_date,
    is_Finished,
    registered_at,
    edited_at
FROM
    review_items
WHERE
    id = sqlc.arg(id)
AND
    user_id = sqlc.arg(user_id);

-- 完了済みの復習日がないか判別するためのクエリ
-- name: HasCompletedReviewDateByItemID :one
SELECT EXISTS (
    SELECT
        1
    FROM
        review_dates
    WHERE
        item_id = sqlc.arg(item_id)
    AND
        user_id = sqlc.arg(user_id)
    AND
        is_completed = TRUE);

-- 復習日Upate処理用。ReviewDateIDを使い回すために使う
-- name: GetReviewDateIDsByItemID :many
SELECT
    id
FROM
    review_dates
WHERE
    item_id = sqlc.arg(item_id)
AND
    user_id = sqlc.arg(user_id)
ORDER BY
    step_number;

-- 移動、完了、学習日変更、その他編集に使う
-- name: UpdateItem :exec
UPDATE
    review_items
SET
    category_id = sqlc.arg(category_id),
    box_id = sqlc.arg(box_id),
    pattern_id = sqlc.arg(pattern_id),
    name = sqlc.arg(name),
    detail = sqlc.arg(detail),
    learned_date = sqlc.arg(learned_date),
    is_Finished = sqlc.arg(is_Finished),
    edited_at = sqlc.arg(edited_at)
WHERE
    id = sqlc.arg(id)
AND
    user_id = sqlc.arg(user_id);

-- 復習日手動変更、完了、学習日変更機能の副次的な変更に使う
-- name: UpdateReviewDates :exec
UPDATE review_dates r
SET
    scheduled_date = v.scheduled_date,
    is_completed = v.is_completed
FROM
    UNNEST(
        sqlc.arg(input)::reviewdate_input[]
    ) AS v(id, scheduled_date, is_completed)
WHERE
    r.id = v.id
AND
    r.user_id = (sqlc.arg(user_id))::uuid;

-- name: UpdateItemAsFinished :exec
UPDATE
    review_items
SET
    is_Finished = true,
    edited_at = sqlc.arg(edited_at)
WHERE
    id = sqlc.arg(id)
AND
    user_id = sqlc.arg(user_id);

-- name: UpdateItemAsUnfinished :exec
UPDATE
    review_items
SET
    is_Finished = false,
    edited_at = sqlc.arg(edited_at)
WHERE
    id = sqlc.arg(id)
AND
    user_id = sqlc.arg(user_id);

-- name: UpdateReviewDateAsCompleted :exec
UPDATE
    review_dates
SET
    is_completed = true
WHERE
    id = sqlc.arg(id)
AND
    user_id = sqlc.arg(user_id);

-- name: UpdateReviewDateAsInCompleted :exec
UPDATE
    review_dates
SET
    is_completed = false
WHERE
    id = sqlc.arg(id)
AND
    user_id = sqlc.arg(user_id);

-- name: GetReviewDatesByItemID :many
SELECT
    id,
    user_id,
    category_id,
    box_id,
    item_id,
    step_number,
    initial_scheduled_date,
    scheduled_date,
    is_completed
FROM
    review_dates
WHERE
    item_id = sqlc.arg(item_id)
AND
    user_id = sqlc.arg(user_id)
ORDER BY
    step_number;

-- name: DeleteItem :exec
DELETE
FROM
    review_items
WHERE
    id = sqlc.arg(id)
AND
    user_id = sqlc.arg(user_id);

-- 復習日のパターンIDがnilに変更されたとき
-- name: DeleteReviewDates :exec
DELETE
FROM
    review_dates
WHERE
    item_id = sqlc.arg(item_id)
AND
    user_id = sqlc.arg(user_id);


-- ボックス内画面用の未完了の全復習物一覧取得機能（復習物（親）のみ一覧取得）
-- name: GetAllUnFinishedItemsByBoxID :many
SELECT
    id,
    user_id,
    category_id,
    box_id,
    pattern_id,
    name,
    detail,
    learned_date,
    is_Finished,
    registered_at,
    edited_at
FROM
    review_items
WHERE
    box_id = sqlc.arg(box_id)
AND
    is_Finished = false
AND
    user_id = sqlc.arg(user_id)
ORDER BY
    registered_at;

--　ボックス内画面用の全復習物一覧取得機能（復習日（子）のみ一覧取得（親は区別しない。親が未完了復習物かどうかも区別しない））。
-- name: GetAllReviewDatesByBoxID :many
SELECT
    id,
    user_id,
    category_id,
    box_id,
    item_id,
    step_number,
    initial_scheduled_date,
    scheduled_date,
    is_completed
FROM
    review_dates
WHERE
    box_id = sqlc.arg(box_id)
AND    
    user_id = sqlc.arg(user_id)
ORDER BY
    item_id,
    step_number;

-- ホーム画面の未分類未完了復習物
-- name: GetAllUnFinishedUnclassifiedItemsByUserID :many
SELECT
    id,
    user_id,
    category_id,
    box_id,
    pattern_id,
    name,
    detail,
    learned_date,
    is_Finished,
    registered_at,
    edited_at
FROM
    review_items
WHERE
    user_id = sqlc.arg(user_id)
AND
    category_id IS NULL
AND
    box_id IS NULL
AND
    is_Finished = false
ORDER BY
    registered_at;

-- name: GetAllUnclassifiedReviewDatesByUserID :many
SELECT
    id,
    user_id,
    category_id,
    box_id,
    item_id,
    step_number,
    initial_scheduled_date,
    scheduled_date,
    is_completed
FROM
    review_dates
WHERE
    user_id = sqlc.arg(user_id)
AND
    category_id IS NULL
AND
    box_id IS NULL
ORDER BY
    item_id,
    step_number;

-- name: GetAllUnFinishedUnclassifiedItemsByCategoryID :many
SELECT
    id,
    user_id,
    category_id,
    box_id,
    pattern_id,
    name,
    detail,
    learned_date,
    is_Finished,
    registered_at,
    edited_at
FROM
    review_items
WHERE
    category_id = sqlc.arg(category_id)
AND
    user_id = sqlc.arg(user_id)
AND
    box_id IS NULL
AND
    is_Finished = false
ORDER BY
    registered_at;

-- name: GetAllUnclassifiedReviewDatesByCategoryID :many
SELECT
    id,
    user_id,
    category_id,
    box_id,
    item_id,
    step_number,
    initial_scheduled_date,
    scheduled_date,
    is_completed
FROM
    review_dates
WHERE
    category_id = sqlc.arg(category_id)
AND
    user_id = sqlc.arg(user_id)
AND
    box_id IS NULL
ORDER BY
    item_id,
    step_number;


-- ここから下は概要表示用の取得クエリ

-- name: CountItemsGroupedByBoxByUserID :many
SELECT
    category_id,
    box_id,
    COUNT(*) AS count
FROM
    review_items
WHERE
    user_id = sqlc.arg(user_id)
AND
    is_Finished = false
GROUP BY
    category_id,
    box_id;


-- name: CountUnclassifiedItemsGroupedByCategoryByUserID :many
SELECT
    category_id,
    COUNT(*) AS count
FROM
    review_items
WHERE
    user_id = sqlc.arg(user_id)
AND
    is_Finished = false
AND
    box_id IS NULL
GROUP BY
    category_id;


-- name: CountUnclassifiedItemsByUserID :many
SELECT
    COUNT(*) AS count
FROM
    review_items
WHERE
    user_id = sqlc.arg(user_id)
AND
    is_Finished = false
AND
    box_id IS NULL;


-- name: CountDailyDatesGroupedByBoxByUserID :many
SELECT
    category_id,
    box_id,
    COUNT(*) AS count
FROM
    review_dates
WHERE
    user_id = sqlc.arg(user_id)
AND
    scheduled_date = sqlc.arg(target_date)
AND
    is_completed = false
AND
    box_id IS NOT NULL
GROUP BY
    category_id,
    box_id;

-- name: CountDailyDatesUnclassifiedGroupedByCategoryByUserID :many
SELECT
    category_id,
    COUNT(*) AS count
FROM
    review_dates
WHERE
    user_id = sqlc.arg(user_id)
AND
    scheduled_date = sqlc.arg(target_date)
AND
    is_completed = false
AND
    box_id IS NULL
GROUP BY
    category_id;

-- name: CountDailyDatesUnclassifiedByUserID :many
SELECT
    COUNT(*) AS count
FROM
    review_dates
WHERE
    user_id = sqlc.arg(user_id)
AND
    scheduled_date = sqlc.arg(target_date)
AND
    is_completed = false
AND
    box_id IS NULL;

-- EditedAt取得専用
-- name: GetEditedAtByItemID :one
SELECT
    edited_at
FROM
    review_items
WHERE
    id = sqlc.arg(id)
AND
    user_id = sqlc.arg(user_id);

-- patternパッケージで使う
-- name: IsPatternRelatedToItemByPatternID :one
SELECT EXISTS (
    SELECT
        1
    FROM
        review_items
    WHERE
        pattern_id = sqlc.arg(pattern_id)
    AND
        user_id = sqlc.arg(user_id)
);

-- 今日の全復習日数を取得
-- name: CountAllDailyReviewDates :one
SELECT
    COUNT(*) AS count
FROM
    review_dates
WHERE
    user_id = sqlc.arg(user_id)
AND
    scheduled_date = sqlc.arg(target_date);

-- LAG→item_idごとにstep_numberの昇順で並べた時、scheduled_dateが持つstep_numberより一個前のstep_numberのscheduled_dateを取得
-- LEAD→item_idごとにstep_numberの昇順で並べた時、scheduled_dateが持つstep_numberより一個後のstep_numberのscheduled_dateを取得
-- 今日の復習日を取得するクエリ
-- name: GetAllDailyReviewDates :many
SELECT
    rd.id,
    rd.category_id,
    rd.box_id,
    rd.step_number,
    rd.prev_scheduled_date,
    rd.scheduled_date,
    rd.next_scheduled_date,
    rd.is_completed,
    ri.id AS item_id,
    ri.name,
    ri.detail,
    ri.learned_date,
    ri.registered_at,
    ri.edited_at
FROM (
    SELECT
        id,
        category_id,
        box_id,
        item_id,
        step_number,
        scheduled_date,
        is_completed,
        CAST(
            LAG(scheduled_date) OVER (
        PARTITION BY item_id
        ORDER BY step_number
        ) AS date
        ) AS prev_scheduled_date,
        CAST(
            LEAD(scheduled_date) OVER (
        PARTITION BY item_id
        ORDER BY step_number
        ) AS date
        ) AS next_scheduled_date
    FROM
        review_dates
    WHERE
        user_id = sqlc.arg(user_id)::uuid
) AS rd
JOIN
    review_items AS ri
ON
    ri.id = rd.item_id
WHERE
    rd.scheduled_date = sqlc.arg(today)::date
ORDER BY
    rd.category_id    NULLS LAST,
    rd.box_id         NULLS LAST,
    ri.registered_at;



-- ボックス内画面用の完了の全復習物一覧取得系（復習物（親）のみ一覧取得）
-- name: GetFinishedItemsByBoxID :many
SELECT
    id,
    user_id,
    category_id,
    box_id,
    pattern_id,
    name,
    detail,
    learned_date,
    is_Finished,
    registered_at,
    edited_at
FROM
    review_items
WHERE
    box_id = sqlc.arg(box_id)
AND
    is_Finished = true
AND
    user_id = sqlc.arg(user_id)
ORDER BY
    registered_at;


-- name: GetUnclassfiedFinishedItemsByCategoryID :many
SELECT
    id,
    user_id,
    category_id,
    box_id,
    pattern_id,
    name,
    detail,
    learned_date,
    is_Finished,
    registered_at,
    edited_at
FROM
    review_items
WHERE
    category_id = sqlc.arg(category_id)
AND
    user_id = sqlc.arg(user_id)
AND
    box_id IS NULL
AND
    is_Finished = true
ORDER BY
    registered_at;

-- name: GetUnclassfiedFinishedItemsByUserID :many
SELECT
    id,
    user_id,
    category_id,
    box_id,
    pattern_id,
    name,
    detail,
    learned_date,
    is_Finished,
    registered_at,
    edited_at
FROM
    review_items
WHERE
    user_id = sqlc.arg(user_id)
AND
    category_id IS NULL
AND
    box_id IS NULL
AND
    is_Finished = true
ORDER BY
    registered_at;