-- name: CreateEntry :one
INSERT INTO entries (
  account_id, amount, transaction_type
) VALUES (
  $1, $2, $3
)
RETURNING *;


-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
LIMIT $1
OFFSET $2;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;

-- name: UpdateEntry :exec
UPDATE entries set amount = $2
WHERE id = $1;