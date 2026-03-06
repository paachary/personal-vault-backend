# go-mongo-project

A RESTful API backend built with Go, Gin, and MongoDB for securely managing personal information, financial documents, and identity proofs.

## Features

-   **User Management** – Register and authenticate users with JWT-based authentication
-   **Identity Proofs** – Store and manage PAN, Aadhaar, Passport, and PRAN details
-   **Financial Documents** – Manage bank account credentials and mutual fund folio information
-   **Account Credentials** – Store IRCTC and email account details
-   **Miscellaneous Details** – Flexible key-value store for additional personal data
-   **Secure Authentication** – All sensitive endpoints are protected by JWT middleware
-   **Multi-Factor Authentication (MFA)** – Time-limited 6-digit MFA codes stored in a dedicated MongoDB collection for step-up verification
-   **Credential Encryption** – Sensitive data (PAN, Aadhaar, bank details, passwords, etc.) is encrypted at rest using AES-256-GCM

## Tech Stack

| Layer            | Technology                                                                  |
| ---------------- | --------------------------------------------------------------------------- |
| Language         | Go (1.25+)                                                                  |
| Web Framework    | [Gin](https://github.com/gin-gonic/gin)                                     |
| Database         | MongoDB (via [mongo-driver v2](https://github.com/mongodb/mongo-go-driver)) |
| Auth             | JWT ([golang-jwt/jwt v5](https://github.com/golang-jwt/jwt))                |
| Config           | [godotenv](https://github.com/joho/godotenv)                                |
| Password Hashing | [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)                     |
| Encryption       | AES-256-GCM (via Go standard library `crypto/aes`, `crypto/cipher`)         |
| CORS             | [gin-contrib/cors](https://github.com/gin-contrib/cors)                     |

## Prerequisites

-   [Go 1.25+](https://golang.org/dl/)
-   [MongoDB](https://www.mongodb.com/try/download/community) (running locally or a remote instance)

## Installation & Setup

1. **Clone the repository**

    ```bash
    git clone https://github.com/paachary/go-mongo-project.git
    cd go-mongo-project
    ```

2. **Install dependencies**

    ```bash
    go mod download
    ```

3. **Configure environment variables**

    Create a `.env` file in the project root (see [Environment Variables](#environment-variables) below).

4. **Run the server**

    ```bash
    go run main.go
    ```

    The API server starts on **port 8080** by default.

## Environment Variables

Create a `.env` file in the project root with the following variables:

```env
# MongoDB connection
MONGO_DB_HOST=localhost
MONGO_DB_PORT=27017
MONGO_DB_NAME=your_database_name
MONGO_COLLECTION_NAME=your_collection_name
MONGO_MFA_COLLECTION_NAME=your_mfa_collection_name

# JWT secret key
ACCESS_TOKEN_SECRET=your_jwt_secret_key

# AES-256-GCM encryption key (must be a 64-character hex string representing 32 bytes)
CREDENTIAL_ENCRYPTION_KEY=your_64_char_hex_key
```

> **Note:** You can generate a valid `CREDENTIAL_ENCRYPTION_KEY` with:
> ```bash
> openssl rand -hex 32
> ```

## API Endpoints

The base URL is `http://localhost:8080`.

### Authentication

| Method | Endpoint         | Auth Required | Description                              |
| ------ | ---------------- | ------------- | ---------------------------------------- |
| POST   | `/register-user` | No            | Register a new user                      |
| POST   | `/login`         | No            | Log in and receive a JWT                 |
| POST   | `/request-mfa`   | No            | Generate and store a 6-digit MFA code    |
| POST   | `/verify-mfa`    | No            | Verify a previously generated MFA code  |

**Request body for `/request-mfa`:** `user_name`, `email_id`

**Request body for `/verify-mfa`:** `user_name`, `email_id`, `code`

All endpoints below require a valid JWT token in the `Authorization` header:

```
Authorization: Bearer <token>
```

---

### User

| Method | Endpoint         | Description                                         |
| ------ | ---------------- | --------------------------------------------------- |
| GET    | `/`              | Get all personal details for the authenticated user |
| PUT    | `/user-details`  | Update user profile (email, mobile, address)        |
| PUT    | `/user-password` | Change user password                                |
| GET    | `/all-users`     | Retrieve all users' data (admin)                    |

---

### PAN Details

| Method | Endpoint | Description             |
| ------ | -------- | ----------------------- |
| POST   | `/pan`   | Add PAN card details    |
| PUT    | `/pan`   | Update PAN card details |
| DELETE | `/pan`   | Delete PAN card details |

**Request body fields:** `pan_number`, `issue_date`, `password`

---

### Aadhaar Details

| Method | Endpoint   | Description                 |
| ------ | ---------- | --------------------------- |
| POST   | `/aadhaar` | Add Aadhaar card details    |
| PUT    | `/aadhaar` | Update Aadhaar card details |
| DELETE | `/aadhaar` | Delete Aadhaar card details |

**Request body fields:** `aadhaar_number`, `issue_date`

---

### PRAN Details

| Method | Endpoint | Description         |
| ------ | -------- | ------------------- |
| POST   | `/pran`  | Add PRAN details    |
| PUT    | `/pran`  | Update PRAN details |
| DELETE | `/pran`  | Delete PRAN details |

**Request body fields:** `pran_number`, `password`

---

### Passport Details

| Method | Endpoint    | Description             |
| ------ | ----------- | ----------------------- |
| POST   | `/passport` | Add passport details    |
| PUT    | `/passport` | Update passport details |
| DELETE | `/passport` | Delete passport details |

**Request body fields:** `passport_number`, `issuer_country`, `issue_date`, `expiry_date`, `user_id`, `password`

---

### IRCTC Details

| Method | Endpoint | Description              |
| ------ | -------- | ------------------------ |
| POST   | `/irctc` | Add IRCTC credentials    |
| PUT    | `/irctc` | Update IRCTC credentials |
| DELETE | `/irctc` | Delete IRCTC credentials |

**Request body fields:** `user_name`, `email_id`, `password`

---

### Email Accounts

| Method | Endpoint | Description          |
| ------ | -------- | -------------------- |
| POST   | `/email` | Add email account    |
| PUT    | `/email` | Update email account |
| DELETE | `/email` | Delete email account |

**Request body fields:** `entity`, `email_id`, `password`

---

### Mutual Funds

| Method | Endpoint | Description                |
| ------ | -------- | -------------------------- |
| POST   | `/mf`    | Add mutual fund details    |
| PUT    | `/mf`    | Update mutual fund details |
| DELETE | `/mf`    | Delete mutual fund details |

**Request body fields:** `fund_name`, `folio_number`, `user_id`, `email_id`, `mobile_number`, `login_password`, `transaction_password`, `mpin`

---

### Bank Details

| Method | Endpoint | Description                 |
| ------ | -------- | --------------------------- |
| POST   | `/bank`  | Add bank account details    |
| PUT    | `/bank`  | Update bank account details |
| DELETE | `/bank`  | Delete bank account details |

**Request body fields:** `bank_name`, `customer_id`, `user_id`, `login_password`, `transaction_password`, `mobile_login_pin`, `mobile_transaction_pin`

> **Note:** A unique `id` (UUID) is auto-generated when adding a bank record. Use this `id` field when updating or deleting.

---

### Miscellaneous Details

| Method | Endpoint | Description                  |
| ------ | -------- | ---------------------------- |
| POST   | `/misc`  | Add miscellaneous details    |
| PUT    | `/misc`  | Update miscellaneous details |
| DELETE | `/misc`  | Delete miscellaneous details |

**Request body fields:** `type_code`, `description`, `key_1`, `val_1`, `key_2`, `val_2`

---

## Data Models

### User

```json
{
    "user_name": "string",
    "is_admin": "string",
    "first_name": "string",
    "last_name": "string",
    "email": "string",
    "mobile_number": "string",
    "password": "string",
    "date_of_birth": "string",
    "address": {
        "street": "string",
        "city": "string",
        "state": "string",
        "postal_code": "string",
        "country": "string"
    },
    "pan_details": [],
    "aadhaar_details": [],
    "pran_details": [],
    "passport_details": [],
    "irctc_details": [],
    "emails": [],
    "mutual_funds": [],
    "bank_details": [],
    "misc_details": []
}
```

### MFA Code

```json
{
    "user_name": "string",
    "email_id": "string",
    "code": "string",
    "createdAt": "timestamp",
    "expiresAt": "timestamp",
    "verified": "boolean"
}
```

> **Note:** MFA codes expire after **5 minutes**. Expired or used codes are automatically deleted.

## Project Structure

```
go-mongo-project/
├── apis/               # HTTP handler functions for each resource
│   ├── aadhaar.go
│   ├── all-details.go
│   ├── bank.go
│   ├── emails.go
│   ├── irctc.go
│   ├── mfa.go
│   ├── misc-details.go
│   ├── mutual-funds.go
│   ├── pan.go
│   ├── passport.go
│   ├── pran.go
│   └── user.go
├── config/             # Application constants and configuration keys
│   └── constants.go
├── db/                 # MongoDB connection and collection initialization
│   ├── collection.go
│   └── db.go
├── middlewares/        # JWT authentication middleware
│   └── auth.go
├── models/             # Data types and MongoDB document operations
│   ├── topFuncs.go
│   └── types.go
├── routes/             # Route registration
│   └── routes.go
├── utils/              # Utility functions (JWT, hashing, UUID, encryption, MFA)
│   ├── encrypt.go
│   ├── hash.go
│   ├── jwt.go
│   ├── mfa.go
│   └── uuid.go
├── go.mod
├── go.sum
└── main.go
```

## CORS Configuration

By default, the server allows requests from `http://localhost:5173` (Vite dev server). To change the allowed origin, update the `config.AllowOrigins` slice in `main.go`.

## License

This project is open source. See the repository for details.
