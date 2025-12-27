-- name: CreateUser :exec
INSERT INTO users (username, master_password_hash, session_token) VALUES ($1, $2, $3);

-- name: CreateVaultItem :exec
INSERT INTO vault_items (vault_id, vault_item_name_encrypted) VALUES ($1, $2);

-- name: CreateCredential :exec
INSERT INTO credentials (vault_item_id, encryptedCredName, encryptedCredPassword) VALUES ($1, $2, $3);

-- name: DeleteUserByID :exec
DELETE FROM users WHERE userid = $1;

-- name: DeleteVaultByID :exec
DELETE FROM vaults WHERE vault_id = $1;

-- name: DeleteVaultItemByID :exec
DELETE FROM vault_items WHERE vault_item_id = $1;

-- name: DeleteCredentialByID :exec
DELETE FROM credentials WHERE credential_id = $1;

-- name: GetUserByUsername :one
SELECT userid, username, master_password_hash, session_token FROM users WHERE username = $1;

-- name: InsertSessionTokenWithUsername :exec
UPDATE users SET session_token = $1 WHERE username = $2;

-- name: GetAllCredentials :many
SELECT c.encryptedCredName, c.encryptedCredPassword, c.credential_id FROM credentials c WHERE vault_item_id = $1;
