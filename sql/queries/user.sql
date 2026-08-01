-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByPhone :one
SELECT *
FROM users
WHERE phone = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (id, phone, created_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: SetUsername :exec
UPDATE users
SET username = $1
WHERE id = $2;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $1
WHERE id = $2;

-- name: ClearUserPassword :exec
UPDATE users
SET password_hash = NULL
WHERE id = $1;

-- name: UpdateUserAvatar :exec
UPDATE users
SET avatar_url = $1
WHERE id = $2;

-- name: UserWithUsernameExists :one
SELECT EXISTS (SELECT 1
               FROM users
               WHERE username = $1);

-- name: UserWithPhoneExists :one
SELECT EXISTS (SELECT 1
               FROM users
               WHERE phone = $1);

-- name: UpsertPhoneOTP :exec
INSERT INTO phone_otp_codes (phone, code_hash, expires_at, attempts, sent_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (phone) DO UPDATE
    SET code_hash  = EXCLUDED.code_hash,
        expires_at = EXCLUDED.expires_at,
        attempts   = EXCLUDED.attempts,
        sent_at    = EXCLUDED.sent_at,
        created_at = EXCLUDED.created_at;

-- name: GetPhoneOTP :one
SELECT *
FROM phone_otp_codes
WHERE phone = $1;

-- name: IncrementPhoneOTPAttempts :exec
UPDATE phone_otp_codes
SET attempts = attempts + 1
WHERE phone = $1;

-- name: DeletePhoneOTP :exec
DELETE
FROM phone_otp_codes
WHERE phone = $1;
