-- name: CreateTransaction :one
INSERT INTO transaction(
    from_account_id, to_account_id, amount
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetTransaction :one
SELECT * from transaction
WHERE id = $1 LIMIT 1;

-- name: GetTransactionBySenderAcc :many
SELECT * from transaction
WHERE from_account_id = $1
LIMIT $2
OFFSET $3;

-- name: GetTransactionByReceiverAcc :many
SELECT * from transaction
WHERE to_account_id = $1
LIMIT $2
OFFSET $3;

-- name: GetTransactionBySenderAndReceiver :many
SELECT * from transaction
WHERE from_account_id = $1 AND to_account_id = $2
LIMIT $3
OFFSET $4;

-- name: DeleteTransaction :exec
DELETE FROM transaction
WHERE id = $1;

-- name: UpdateTransaction :exec
UPDATE transaction set amount = $2
WHERE id = $1;