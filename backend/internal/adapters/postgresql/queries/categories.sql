-- name: CreateCategory :one
INSERT INTO categories (
    id,
    name
) VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetCategoryByID :one
SELECT *
FROM categories
WHERE id = $1;

-- name: ListCategories :many
SELECT *
FROM categories
ORDER BY name ASC;

-- name: UpdateCategoryName :one
UPDATE categories
SET
    name = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :execrows
DELETE FROM categories
WHERE id = $1;
