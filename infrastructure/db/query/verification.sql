-- name: CreateEmailVerification :exec
INSERT INTO email_verifications (
    id,
    user_id,
    code_hash,
    expires_at
) VALUES (
    sqlc.arg(id),
    sqlc.arg(user_id),
    sqlc.arg(code_hash),
    sqlc.arg(expires_at)
);

-- name: FindEmailVerificationByUserID :one
SELECT
    id,
    user_id,
    code_hash,
    expires_at
FROM
    email_verifications
WHERE
    user_id = sqlc.arg(user_id)
LIMIT 1;

-- name: DeleteEmailVerification :exec
DELETE FROM
    email_verifications
WHERE
    id = sqlc.arg(id);

-- name: DeleteEmailVerificationByUserID :exec
DELETE FROM
    email_verifications
WHERE
    user_id = sqlc.arg(user_id);