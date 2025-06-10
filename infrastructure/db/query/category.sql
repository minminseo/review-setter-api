-- name: CreateCategory :exec
INSERT INTO 
    categories (
        id,
        user_id,
        name,
        registered_at,
        edited_at
    ) VALUES (
        sqlc.arg(id),
        sqlc.arg(user_id),
        sqlc.arg(name),
        sqlc.arg(registered_at),
        sqlc.arg(edited_at)
    );

-- name: GetAllCategoriesByUserID :many
SELECT
    id,
    user_id,
    name,
    registered_at,
    edited_at
FROM
    categories
WHERE
    user_id = sqlc.arg(user_id)
ORDER BY
    registered_at;

-- name: GetCategoryByID :one
SELECT
    name,
    registered_at,
    edited_at
FROM
    categories
WHERE
    id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);

-- name: UpdateCategory :exec
UPDATE
    categories
SET
    name = sqlc.arg(name),
    edited_at = sqlc.arg(edited_at)
WHERE
    id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);

-- name: DeleteCategory :exec
DELETE 
FROM
    categories
WHERE
    id = sqlc.arg(id) AND user_id = sqlc.arg(user_id);

-- item_usecaseで使うクエリ
-- name: GetCategoryNamesByCategoryIDs :many
-- args: category_ids uuid[]
SELECT
    id,
    name
FROM
    categories
WHERE
    id = ANY(sqlc.arg(category_ids)::uuid[]);

