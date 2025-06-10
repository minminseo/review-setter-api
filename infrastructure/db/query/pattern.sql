-- name: CreatePattern :exec
INSERT INTO
    review_patterns (
        id,
        user_id,
        name,
        target_weight,
        registered_at,
        edited_at
    )
VALUES (
        sqlc.arg(id),
        sqlc.arg(user_id),
        sqlc.arg(name),
        sqlc.arg(target_weight),
        sqlc.arg(registered_at),
        sqlc.arg(edited_at)
    );

-- 新規一括挿入時と、一括更新時に使う
-- name: CreatePatternSteps :copyfrom
INSERT INTO
    pattern_steps (
        id,
        user_id,
        pattern_id,
        step_number,
        interval_days
    ) VALUES (
        sqlc.arg(id),
        sqlc.arg(user_id),
        sqlc.arg(pattern_id),
        sqlc.arg(step_number),
        sqlc.arg(interval_days)
    );


-- 復習パターンそのものが更新対象かどうか判定するために使う
-- name: GetPatternByID :one
SELECT
    id,
    user_id,
    name,
    target_weight,
    registered_at,
    edited_at
FROM
    review_patterns
WHERE
    id = sqlc.arg(id)
AND
    user_id = sqlc.arg(user_id);

-- 復習ステップが更新対象かどうか判定するために使う
-- name: GetPatternStepsByPatternID :many
SELECT
    id,
    user_id,
    pattern_id,
    step_number,
    interval_days
FROM
    pattern_steps
WHERE
    pattern_id = sqlc.arg(pattern_id)
AND
    user_id = sqlc.arg(user_id)
ORDER BY
    step_number;

-- pattern系のリクエストで、更新対象の中に復習パターンそのものが含まれる場合に発行するクエリ
-- name: UpdatePattern :exec
UPDATE
    review_patterns
SET
    name = sqlc.arg(name),
    target_weight = sqlc.arg(target_weight),
    edited_at = sqlc.arg(edited_at)
WHERE
    id = sqlc.arg(id)
AND
    user_id = sqlc.arg(user_id);



-- name: DeletePattern :exec
DELETE
FROM
    review_patterns
WHERE
    id = sqlc.arg(id)
AND
    user_id = sqlc.arg(user_id);

-- 復習ステップが更新対象に含まれた場合に発行する一括削除用のクエリ
-- name: DeletePatternSteps :exec
DELETE
FROM
    pattern_steps
WHERE
    pattern_id = sqlc.arg(pattern_id)
AND
    user_id = sqlc.arg(user_id);


-- 全パターン取得機能（パターン（親）のみ一覧取得）
-- name: GetAllPatternsByUserID :many
SELECT
    id,
    user_id,
    name,
    target_weight,
    registered_at,
    edited_at
FROM
    review_patterns
WHERE
    user_id = sqlc.arg(user_id)
ORDER BY
    registered_at;

--　全パターン取得機能（ステップ（子）のみ一覧取得（親は区別しない））
-- name: GetAllPatternStepsByUserID :many
SELECT
    id,
    user_id,
    pattern_id,
    step_number,
    interval_days
FROM
    pattern_steps
WHERE
    user_id = sqlc.arg(user_id)
ORDER BY
    pattern_id,
    step_number;

-- item_usecaseで使うクエリ。
-- name: GetPatternTargetWeightsByPatternIDs :many
-- args: pattern_ids uuid[]
SELECT
    id,
    target_weight
FROM
    review_patterns
WHERE
    id = ANY(sqlc.arg(pattern_ids)::uuid[]);
