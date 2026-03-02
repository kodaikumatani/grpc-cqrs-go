-- name: CreateRecipe :exec
INSERT INTO recipes (id, user_id, title, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetRecipe :one
SELECT id, user_id, title, description, created_at, updated_at
FROM recipes
WHERE id = $1;
