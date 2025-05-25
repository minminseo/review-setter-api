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

--name: FindUserByEmail: one
SELECT
    id,
    email,
    password,
    timezone,
    theme_color,
    language,
FROM
    users
WHERE
    email = sqlc.arg(email);

--name: UpdateUser: exec
UPDATE
    users
SET
    email = sqlc.arg(email),
    password = sqlc.arg(password),
    timezone = sqlc.arg(timezone),
    theme_color = sqlc.arg(theme_color),
    language = sqlc.arg(language),
WHERE
    id = sqlc.arg(id);