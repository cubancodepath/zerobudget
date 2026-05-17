-- name: CreatePayee :one
INSERT INTO payees (
    id,
    name
) VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetPayeeByID :one
SELECT *
FROM payees
WHERE id = $1;

-- name: UpsertPayeeByNameCaseInsensitive :one
INSERT INTO payees (
    id,
    name
) VALUES (
    $1,
    $2
)
ON CONFLICT (name) DO UPDATE
SET updated_at = now()
RETURNING *;

-- name: ListPayees :many
SELECT *
FROM payees
ORDER BY name ASC;

-- name: UpdatePayeeName :one
UPDATE payees
SET
    name = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePayee :execrows
DELETE FROM payees
WHERE id = $1;
