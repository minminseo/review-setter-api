-- name: CreateUser :exec
INSERT INTO 
    users (
        id,
        email,
        password,
        timezone,
        theme_color,
        language
    ) VALUES (
        sqlc.arg(id),
        sqlc.arg(email),
        sqlc.arg(password),
        sqlc.arg(timezone),
        sqlc.arg(theme_color),
        sqlc.arg(language)
    );

-- name: FindUserByEmail :one
SELECT
    id,
    password,
    theme_color,
    language,
    verified_at
FROM
    users
WHERE
    email = sqlc.arg(email);

-- name: GetUserSettingByID :one
SELECT
    email,
    timezone,
    theme_color,
    language
FROM
    users
WHERE
    id = sqlc.arg(id);

-- name: UpdateUser :exec
UPDATE
    users
SET
    email = sqlc.arg(email),
    timezone = sqlc.arg(timezone),
    theme_color = sqlc.arg(theme_color),
    language = sqlc.arg(language)
WHERE
    id = sqlc.arg(id);

-- name: UpdateUserPassword :exec
UPDATE
    users
SET
    password = sqlc.arg(password)
WHERE
    id = sqlc.arg(id);

-- name: UpdateVerifiedAt :exec
UPDATE
    users
SET
    verified_at = sqlc.arg(verified_at)
WHERE
    id = sqlc.arg(id);