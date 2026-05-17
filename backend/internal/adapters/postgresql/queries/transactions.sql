-- name: CreateTransaction :one
INSERT INTO transactions (
    id,
    account_id,
    category_id,
    payee_id,
    amount_cents,
    transaction_date,
    is_reconciled,
    note
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetTransactionByID :one
SELECT *
FROM transactions
WHERE id = $1;

-- name: ListTransactionsByAccount :many
SELECT *
FROM transactions
WHERE account_id = $1
ORDER BY transaction_date DESC, created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListTransactionsByAccountAndDateRange :many
SELECT *
FROM transactions
WHERE account_id = $1
  AND transaction_date >= $2
  AND transaction_date <= $3
ORDER BY transaction_date DESC, created_at DESC
LIMIT $4
OFFSET $5;

-- name: UpdateTransaction :one
UPDATE transactions
SET
    account_id = $2,
    category_id = $3,
    payee_id = $4,
    amount_cents = $5,
    transaction_date = $6,
    is_reconciled = $7,
    note = $8,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: SetTransactionReconciled :one
UPDATE transactions
SET
    is_reconciled = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTransaction :execrows
DELETE FROM transactions
WHERE id = $1;
