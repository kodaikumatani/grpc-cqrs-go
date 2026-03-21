-- name: CreateRecipe :exec
INSERT INTO recipes (id, user_id, title, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetRecipe :one
SELECT id, user_id, title, description, created_at, updated_at
FROM recipes
WHERE id = $1;

-- name: UpdateRecipe :exec
UPDATE recipes
SET title = $2, description = $3, updated_at = $4
WHERE id = $1;

-- name: GetRecipeWithUser :one
SELECT r.id, r.user_id, r.title, r.description, r.created_at, r.updated_at,
       u.name AS user_name, u.email AS user_email
FROM recipes r
JOIN users u ON r.user_id = u.id
WHERE r.id = $1;
