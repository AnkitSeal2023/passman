CREATE TABLE users (
    userid INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    master_password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE vaults(
    vault_id INTEGER PRIMARY KEY,
    userid INTEGER NOT NULL,
    vault_item_name_encrypted VARCHAR(255) UNIQUE NOT NULL,
    FOREIGN KEY (userid) REFERENCES users(userid) ON DELETE CASCADE
);

CREATE TABLE vault_items (
    vault_item_id INTEGER PRIMARY KEY,
    vault_id INTEGER NOT NULL,
    vault_item_name_encrypted VARCHAR(255) UNIQUE NOT NULL,
    FOREIGN KEY (vault_id) REFERENCES vaults(vault_id) ON DELETE CASCADE
);

CREATE TABLE credentials (
    credential_id INTEGER PRIMARY KEY,
    vault_item_id INTEGER NOT NULL,
    encryptedCredName VARCHAR(255) NOT NULL,
    encryptedCredPassword VARCHAR(255) NOT NULL,
    FOREIGN KEY (vault_item_id) REFERENCES vault_items(vault_item_id) ON DELETE CASCADE
);
