-- name: CreateUser :exec
INSERT INTO users (username, master_password_hash, session_token, encr_dek) VALUES ($1, $2, $3, $4);

-- name: CreateVault :exec
INSERT INTO vaults (vault_id, userid, vault_item_name_encrypted) VALUES ($1, $2, $3);

-- name: GetVaultByUserId :one
SELECT vault_id, userid, vault_item_name_encrypted FROM vaults WHERE userid = $1;

-- name: CreateVaultItem :exec
INSERT INTO vault_items (vault_item_id, vault_id, vault_item_name_encrypted) VALUES ($1, $2, $3);

-- name: GetVaultItemByNameAndVaultId :one
SELECT vault_item_id, vault_id, vault_item_name_encrypted FROM vault_items WHERE vault_id = $1 AND vault_item_name_encrypted = $2;

-- name: GetMaxVaultId :one
SELECT COALESCE(MAX(vault_id), 0) FROM vaults;

-- name: GetMaxVaultItemId :one
SELECT COALESCE(MAX(vault_item_id), 0) FROM vault_items;

-- name: GetMaxCredentialId :one
SELECT COALESCE(MAX(credential_id), 0) FROM credentials;

-- name: CreateCredential :exec
INSERT INTO credentials (credential_id, vault_item_id, encryptedCredName, encryptedCredPassword) VALUES ($1, $2, $3, $4);

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

-- name: GetVaultItemsByUserAndName :many
SELECT vi.vault_item_id, vi.vault_item_name_encrypted FROM vault_items vi
JOIN vaults v ON vi.vault_id = v.vault_id
JOIN users u ON v.userid = u.userid
WHERE u.username = $1 AND vi.vault_item_name_encrypted = $2;

-- name: GetAllVaultItemsByUser :many
SELECT DISTINCT vi.vault_item_name_encrypted FROM vault_items vi
JOIN vaults v ON vi.vault_id = v.vault_id
JOIN users u ON v.userid = u.userid
WHERE u.username = $1;

-- name: GetVaultItemsForUser :many
SELECT vi.vault_item_id, vi.vault_item_name_encrypted FROM vault_items vi
JOIN vaults v ON vi.vault_id = v.vault_id
JOIN users u ON v.userid = u.userid
WHERE u.username = $1;

-- name: GetUserDEK :one
SELECT encr_dek FROM users WHERE username = $1;
