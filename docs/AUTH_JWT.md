# JWT Authentication Implementation Plan

**Objective**: Implement complete JWT-based authentication for MicroSQL-AGo backend, enabling secure token generation on login and validation for protected endpoints.

---

## Current State

- **Existing components**:
  - User model: `internal/domain/entities/user.go` (ID, Username, Email, Password, Role, CreatedAt, LastLogin, IsActive)
  - User handler: `internal/adapters/primary/http/handlers/user_handler.go` (Register, Login methods)
  - Config: JWT_SECRET stored in `internal/config/config.go`
  - HTTP middleware directory: `internal/adapters/primary/http/middleware/` exists with auth_middleware.go
  
- **Missing/incomplete**:
  - No JWT token generation logic in login
  - No JWT verification middleware
  - No token claims structure defined
  - No token expiry/refresh token support
  - No protected endpoint middleware integration

---

## Implementation Plan

### Phase 1: JWT Token Generation (Login)
**Files to create/modify**:
- `internal/adapters/secondary/security/jwt.go` — JWT token generation and parsing utilities
- `internal/adapters/primary/http/handlers/user_handler.go` — Update Login() to generate and return JWT token

**Steps**:
1. Create `jwt.go` with:
   - `GenerateToken(userID uint, username string, role string, jwtSecret string, expiryHours int) (string, error)`
   - `ValidateToken(tokenString string, jwtSecret string) (*TokenClaims, error)` 
   - Define `TokenClaims` struct with `UserID`, `Username`, `Role`, `exp`, `iat`

2. Update `user_handler.Login()` to:
   - Verify credentials (username/password)
   - Generate JWT token using `GenerateToken()`
   - Return response: `{ "token": "<jwt>", "user": { "id": ..., "username": ..., "role": ... } }`
   - Set 24-hour token expiry

3. Update login response DTOs:
   - Create `LoginResponseDTO` with token and user fields

**Expected Result**: Client can call `POST /api/auth/login` with credentials and receive JWT token.

---

### Phase 2: JWT Verification Middleware
**Files to create/modify**:
- `internal/adapters/primary/http/middleware/auth_middleware.go` — JWT validation middleware
- `internal/adapters/primary/http/routes.go` — Register middleware on protected routes

**Steps**:
1. Create `AuthMiddleware()` function that:
   - Extracts `Authorization: Bearer <token>` header
   - Calls `jwt.ValidateToken()` to parse and validate
   - On success: sets `userID` and `role` in Gin context
   - On failure: returns 401 Unauthorized with error message

2. Create `RoleMiddleware(allowedRoles ...string)` for role-based access:
   - Checks if user's role is in allowedRoles
   - Returns 403 Forbidden if not authorized

3. Register middleware on protected routes:
   - `/api/audits/execute` — requires auth (role: any)
   - Other protected endpoints as needed

**Expected Result**: Protected routes validate JWT token and extract user identity.

---

### Phase 3: Integration and Testing
**Files to create/modify**:
- `internal/adapters/primary/http/routes.go` — Register middleware on routes
- `internal/domain/usecases/controls/execute_audit.go` — Already uses userID from context (no changes needed, just verify)
- Tests: `internal/domain/usecases/.../*_test.go` (optional but recommended)

**Steps**:
1. Apply auth middleware to protected routes in `RegisterRoutes()`
2. Test login endpoint with curl/Postman:
   ```bash
   POST /api/auth/login
   Content-Type: application/json
   
   { "username": "user1", "password": "pass123" }
   
   Response:
   { "token": "<jwt>", "user": { "id": 1, "username": "user1", "role": "admin" } }
   ```

3. Test protected endpoint (audit) with token:
   ```bash
   POST /api/audits/execute
   Authorization: Bearer <token>
   Content-Type: application/json
   
   { "control_ids": [1, 2], "database": "auditeddb" }
   ```

4. Test invalid token → 401 response

---

## Implementation Steps (Execution Order)

### Step 1: JWT Utilities (`internal/adapters/secondary/security/jwt.go`)
- [ ] Define `TokenClaims` struct
- [ ] Implement `GenerateToken()`
- [ ] Implement `ValidateToken()`
- [ ] Add error handling for expired/invalid tokens

### Step 2: Update Login Handler (`internal/adapters/primary/http/handlers/user_handler.go`)
- [ ] Create login request/response DTOs
- [ ] Update `Login()` to call `GenerateToken()`
- [ ] Return token in response
- [ ] Handle password validation and user lookup

### Step 3: Middleware (`internal/adapters/primary/http/middleware/auth_middleware.go`)
- [ ] Implement `AuthMiddleware()` to extract and validate token
- [ ] Implement `RoleMiddleware()` for role checks
- [ ] Set userID and role in Gin context

### Step 4: Route Registration (`internal/adapters/primary/http/routes.go`)
- [ ] Register auth middleware on protected routes
- [ ] Ensure `/api/auth/login` is unprotected
- [ ] Ensure `/api/audits/execute` requires auth

### Step 5: Testing
- [ ] Manual test: login endpoint
- [ ] Manual test: protected endpoint with valid token
- [ ] Manual test: protected endpoint with invalid/missing token
- [ ] Verify token claims are accessible in handlers

---

## Technical Details

### Token Claim Structure
```go
type TokenClaims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}
```

### Token Expiry
- Default: 24 hours
- Can be configured via env: `JWT_EXPIRY_HOURS` (optional, default 24)

### Error Responses
- **401 Unauthorized**: Missing/invalid token, expired token
- **403 Forbidden**: Valid token but insufficient role permissions
- **400 Bad Request**: Invalid login credentials

### Dependencies
- `github.com/golang-jwt/jwt/v5` (likely already in go.mod; verify with `go list`)

---

## Files to Be Created/Modified

| File | Action | Purpose |
|------|--------|---------|
| `internal/adapters/secondary/security/jwt.go` | CREATE | JWT token generation/validation |
| `internal/adapters/primary/http/handlers/user_handler.go` | MODIFY | Update Login() to return JWT |
| `internal/adapters/primary/http/handlers/user_handler.go` | ADD | LoginResponseDTO |
| `internal/adapters/primary/http/middleware/auth_middleware.go` | MODIFY | Add JWT validation |
| `internal/adapters/primary/http/routes.go` | MODIFY | Register middleware on routes |
| `internal/config/config.go` | MODIFY | Add JWT_EXPIRY_HOURS (optional) |

---

## Success Criteria

✅ Client can login and receive JWT token  
✅ Token is valid for 24 hours  
✅ Protected endpoints require valid token  
✅ Invalid/expired tokens return 401  
✅ Insufficient role permissions return 403  
✅ userID and role are available in handlers via Gin context  

---

## Timeline Estimate

- **Phase 1 (JWT utils + Login)**: ~30-45 min
- **Phase 2 (Middleware)**: ~20-30 min  
- **Phase 3 (Integration + Testing)**: ~15-20 min
- **Total**: ~1.5-2 hours

---

## Notes

- Token validation happens in middleware, executed before each request to protected endpoints.
- User credentials are verified using password comparison (ensure hashing is consistent).
- Role-based access can be extended to support granular permissions in future.
- Consider adding token refresh endpoint for long-lived sessions (outside current scope).
