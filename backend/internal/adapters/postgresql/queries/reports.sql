-- name: GetAccountBalance :one
SELECT
    a.id AS account_id,
    a.initial_balance_cents + COALESCE(SUM(t.amount_cents), 0)::BIGINT AS balance_cents
FROM accounts a
LEFT JOIN transactions t ON t.account_id = a.id
WHERE a.id = $1
GROUP BY a.id, a.initial_balance_cents;

-- name: GetReconciledAccountBalance :one
SELECT
    a.id AS account_id,
    a.initial_balance_cents + COALESCE(SUM(t.amount_cents), 0)::BIGINT AS balance_cents
FROM accounts a
LEFT JOIN transactions t
    ON t.account_id = a.id
   AND t.is_reconciled = true
WHERE a.id = $1
GROUP BY a.id, a.initial_balance_cents;
