-- name: CreateCategory :exec
INSERT INTO 
    categories (
        id,
        user_id,
        name
    ) VALUES (
        sqlc.arg(id),
        sqlc.arg(user_id),
        sqlc.arg(name)
    );

-- name: GetAllCategoriesByUserID :many
SELECT
    id,
    user_id,
    name
FROM
    categories
WHERE
    user_id = sqlc.arg(user_id)
ORDER BY
    created_at;

-- name: GetCategoryByID :one
SELECT
    name
FROM
    categories
WHERE
    id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);

-- name: UpdateCategory :exec
UPDATE
    categories
SET
    name = sqlc.arg(name)
WHERE
    id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);

-- name: DeleteCategory :exec
DELETE FROM
    categories
WHERE
    id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);