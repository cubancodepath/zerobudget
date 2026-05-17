-- name: CreateAccount :one
INSERT INTO accounts (
    id,
    name,
    type,
    currency_code,
    initial_balance_cents,
    is_active
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetAccountByID :one
SELECT *
FROM accounts
WHERE id = $1;

-- name: ListActiveAccounts :many
SELECT *
FROM accounts
WHERE is_active = true
ORDER BY created_at DESC;

-- name: UpdateAccount :one
UPDATE accounts
SET
    name = $2,
    type = $3,
    currency_code = $4,
    initial_balance_cents = $5,
    is_active = $6,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeactivateAccount :exec
UPDATE accounts
SET
    is_active = false,
    updated_at = now()
WHERE id = $1;
