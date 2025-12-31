# Passman

A secure password manager web application built with Go.

## Features

- 🔐 Secure password storage with encryption
- 🔑 Master password authentication
- 🗄️ PostgreSQL database backend
- 🐳 Docker support for easy deployment
- 🛡️ Argon2id password hashing
- 🔒 AES-256 encryption for credentials

## Tech Stack

- **Backend:** Go 1.24.5
- **Database:** PostgreSQL
- **Templating:** templ
- **Encryption:** AES-256 GCM with Argon2id KDF
- **Containerization:** Docker

## Quick Start

### Using Docker

```bash
docker build -t passman .
docker run -p 5000:5000 passman
```

### Local Development

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/api/main.go
```

Access the application at `http://localhost:5000`

## Environment Variables

- `PORT` - Server port (default: 5000)
- Database connection string should be configured in the application

## Security

- Master passwords are hashed using Argon2id
- Credentials are encrypted with AES-256-GCM
- Data Encryption Keys (DEK) are derived from master passwords
- Session-based authentication

## License

MIT