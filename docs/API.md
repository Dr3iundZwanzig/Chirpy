# Chirpy API Documentation

## Overview

Chirpy is a simple microblogging API for creating and fetching short messages called "chirps." This document describes available endpoints, authentication, models, examples, and common error responses.

## Base URL

- `http://localhost:8080` (default for local development)

## Authentication

- Chirpy uses JWT access tokens. Obtain a token via `POST /login` and include it in requests as an `Authorization: Bearer <token>` header.
- Refresh tokens are supported via `POST /refresh`.

## Endpoints

- **Auth**
  - `POST /login` — Exchange email/password for access and refresh tokens.
    - Body: `{ "email": "user@example.com", "password": "secret" }`
    - Response: `{ "access_token": "...", "refresh_token": "..." }`

  - `POST /refresh` — Exchange a refresh token for a new access token.
    - Body: `{ "refresh_token": "..." }`

- **Users**
  - `GET /users` — List users (may be paginated).
  - `GET /users/{id}` — Get user by UUID.
  - `POST /users` — Create a new user.
    - Body example: `{ "display_name": "Alice", "email": "alice@example.com", "password": "secret" }`

- **Chirps**
  - `GET /chirps` — List chirps across all users. Query params: `limit`, `offset`, `order`.
  - `GET /chirps/{id}` — Get a chirp by ID.
  - `GET /users/{id}/chirps` — List chirps for a specific user (ordered by `created_at`).
  - `POST /chirps` — Create a new chirp (requires auth).
    - Body example: `{ "body": "Hello world!" }`

- **Password Reset & Webhooks**
  - `POST /reset` — Request password reset (email-based flow).
  - Webhook endpoints are available under `/webhooks` for external integrations (see implementation-specific docs).

## Models

- **User**
  - `id`: UUID
  - `display_name`: string
  - `email`: string
  - `created_at`: timestamp

- **Chirp**
  - `id`: integer
  - `user_id`: UUID
  - `body`: string
  - `created_at`: timestamp

- **RefreshToken**
  - `id`: UUID
  - `user_id`: UUID
  - `expires_at`: timestamp

## Error Responses

- Standard JSON error: `{ "error": "message" }` with appropriate HTTP status codes:
  - `400` — Bad request / validation error
  - `401` — Unauthorized (missing/invalid token)
  - `403` — Forbidden
  - `404` — Not found
  - `500` — Internal server error

## Examples

- Obtain token (curl):

```
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret"}'
```

- Create chirp (curl):

```
curl -X POST http://localhost:8080/chirps \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"body":"Hello from API"}'
```

- Get user chirps (curl):

```
curl http://localhost:8080/users/<user-uuid>/chirps
```

## Notes and Implementation Details

- The SQL query `SELECT * FROM chirps WHERE user_id = $1 ORDER BY created_at ASC` will only match rows with exactly the provided `user_id`. Passing an all-zero UUID (e.g. `00000000-0000-0000-0000-000000000000`) will not return all rows — it matches only rows whose `user_id` equals that zero value.
- Pagination and rate-limiting are recommended for production use.

## Where to look in the code

- Handlers: `handler*.go` files implement routes and request handling.
- Models and DB access: `internal/database` contains generated SQL bindings and models.

---
