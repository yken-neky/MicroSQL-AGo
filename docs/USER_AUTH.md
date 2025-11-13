# User Authentication: Register & Login

Fecha: 2025-11-13

Este documento describe la implementación del endpoint de registro y login con JWT para el servicio backend.

## Endpoints

- POST /api/users/register
  - Payload: RegisterRequest
  - Response: LoginResponse (token + user)
  - Status: 201 Created on success

- POST /api/auth/login
  - Payload: LoginRequest
  - Response: LoginResponse (token + user)
  - Status: 200 OK on success

## DTOs

RegisterRequest:
- username (string, required)
- first_name (string, required)
- last_name (string, required)
- email (string, required, email)
- password (string, required, min 8)

LoginRequest:
- username
- password

LoginResponse:
- token: JWT
- user: UserResponse (id, username, email, first_name, last_name, role)

## Register Flow (server-side)

1. Bind and validate JSON payload.
2. Check uniqueness of `username` and `email`. If conflict, return 409.
3. Hash password with `bcrypt.GenerateFromPassword(..., bcrypt.DefaultCost)`.
4. Create `users` record in DB via GORM.
5. Generate JWT token via `JWTService.GenerateToken(user.ID, user.Username, user.Role)`.
6. Set `last_login` to current timestamp (optional).
7. Return 201 with `LoginResponse` containing token + user.

## Login Flow (server-side)

1. Bind and validate JSON payload.
2. Lookup user by username.
3. Compare password with `bcrypt.CompareHashAndPassword`.
4. If valid and user active, generate JWT token via `JWTService`.
5. Update `last_login` timestamp.
6. Return 200 with `LoginResponse` containing token + user.

## Error Cases

- 400 Bad Request: invalid JSON or validation errors.
- 401 Unauthorized: invalid credentials or inactive account.
- 409 Conflict: username or email already exists on registration.
- 500 Internal Server Error: DB errors, hashing failures, or JWT service unavailability.

## Security Notes

- Passwords are hashed using bcrypt before storage.
- JWT tokens are signed with `HS256` using `JWT_SECRET` from configuration.
- Token expiry default: 24 hours.
- Use HTTPS in production to protect tokens in transit.

## Example curl flows

Register:

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"jdoe","first_name":"John","last_name":"Doe","email":"jdoe@example.com","password":"S3curePassw0rd"}'
```

Login:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"jdoe","password":"S3curePassw0rd"}'
```

Access protected endpoint (example):

```bash
curl -X POST http://localhost:8080/api/audits/execute \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"script_id": 1}'
```

## Notes

- The `Register` handler currently creates the user and returns a token immediately (convenience flow). If email verification is desired, change flow to create inactive user and require verification before issuing token.
- Role default: `cliente` (set by code if DB default not applied). Adjust as needed.

