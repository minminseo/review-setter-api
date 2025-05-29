-- name: CreateBox :exec
INSERT INTO 
    review_boxes (
        id,
        user_id,
        category_id,
        pattern_id,
        name,
        registered_at,
        edited_at
    ) VALUES (
        sqlc.arg(id),
        sqlc.arg(user_id),
        sqlc.arg(category_id),
        sqlc.arg(pattern_id),
        sqlc.arg(name),
        sqlc.arg(registered_at),
        sqlc.arg(edited_at)
    );

-- name: GetAllBoxesByCategoryID :many
SELECT
    id,
    user_id,
    category_id,
    pattern_id,
    name,
    registered_at,
    edited_at
FROM
    review_boxes
WHERE
    category_id = sqlc.arg(category_id) AND user_id = sqlc.arg(user_id)
ORDER BY
    registered_at;

-- name: GetBoxByID :one
SELECT
    id,
    user_id,
    category_id,
    pattern_id,
    name,
    registered_at,
    edited_at
FROM
    review_boxes
WHERE
    id = sqlc.arg(id) AND category_id = sqlc.arg(category_id) AND user_id = sqlc.arg(user_id);

-- name: UpdateBox :exec
UPDATE
    review_boxes
SET
    name = sqlc.arg(name),
    edited_at = sqlc.arg(edited_at)
WHERE
    id = sqlc.arg(id) AND category_id = sqlc.arg(category_id) AND user_id = sqlc.arg(user_id);

-- name: UpdateBoxIfNoReviewItems :execrows
UPDATE
    review_boxes
SET
    pattern_id = sqlc.arg(pattern_id),
    name       = sqlc.arg(name),
    edited_at  = sqlc.arg(edited_at)
WHERE
    review_boxes.id = sqlc.arg(id)
AND
    review_boxes.category_id = sqlc.arg(category_id)
AND
    review_boxes.user_id = sqlc.arg(user_id)
AND NOT EXISTS (
            SELECT 
                1
            FROM 
                review_items
            WHERE 
                review_items.box_id  = sqlc.arg(box_id)
            AND 
                review_items.user_id = sqlc.arg(user_id)
            );


-- name: DeleteBox :exec
DELETE FROM
    review_boxes
WHERE
    id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);

-- -- name: CountGroupedByCategoryByUserID :many
-- SELECT
--     c.id AS category_id,
--     c.name AS category_name,
--     COUNT(b.id) AS box_count
-- FROM
--     review_boxes b
-- JOIN
--     categories c ON b.category_id = c.id
-- WHERE
--     b.user_id = sqlc.arg(user_id)
-- GROUP BY
--     c.id, c.name
-- ORDER BY
--     c.registered_at;

