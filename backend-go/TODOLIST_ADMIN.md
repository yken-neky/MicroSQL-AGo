# Admin TODOs and recommendations

This document lists security, implementation and operational recommendations for administering users, sessions, and admin-level operations.

Goals
- Allow safe promotion/demotion of users to/from admin role
- Provide administrative visibility over active sessions (already implemented)
- Provide safe forcible logout of sessions by admin
- Make role changes immediate and auditable
- Protect admin endpoints strictly and limit token exposure

Quick answers
- A user is considered an admin when `users.role == "admin"` in the database.
- The system recognizes an admin via the `role` claim on the user's JWT — the middleware checks the role claim and enforces access via `RequireRole("admin")`.

Important security caveats
- Tokens already issued will keep the previous claims until they expire or are invalidated. To make role changes immediate you must invalidate existing sessions (see below).
- Storing full JWT tokens in DB is convenient for session invalidation and admin listing, but increases attack surface if DB read access is compromised. Consider storing a cryptographic hash of the token and only showing admin a token preview or metadata (e.g., first/last 4 chars, IP, user agent) instead of the full token in production.

Recommended endpoints and flows (high level)
1) Admin: Promote a user to admin
   - POST /api/admin/users/:id/promote
   - Body: optional {"role":"admin"} or default to admin
   - Behavior: validate caller is admin; update `users.role` to `admin`; invalidate user's existing sessions; write an audit log entry stating who promoted whom and why.

2) Admin: Demote user / remove admin privileges
   - POST /api/admin/users/:id/demote
   - Behavior: same as promote but set role to `cliente` (or other default); invalidate sessions; audit.

3) Admin: List & inspect sessions (already implemented)
   - GET /api/admin/sessions
   - Must be restricted to `admin` role (middleware `RequireRole("admin")`)
   - Consider returning only: session_id, user_id, username, email, partial token, expires_at, created_at, last_seen, ip, user_agent

4) Admin: Forcibly revoke a single session
   - POST /api/admin/sessions/:session_id/revoke
   - Behavior: set `sessions.is_active = false` for the session; optionally add reason and who performed the revoke to an audit table; notify user if needed.

5) Session model and token storage (best practices)
   - Prefer storing a one-way hash (HMAC or bcrypt-like) of the token for admin checks / invalidation matching — store salt and algorithm.
   - If you must store tokens in plain-text for operational reasons, restrict DB access and encrypt at rest; avoid exposing tokens in admin UIs.
   - Use short-lived access tokens (e.g., 15–60 minutes) + refresh tokens so role changes propagate quickly.

6) Token invalidation & immediate role-change semantics
   - When promoting/demoting a user, also invalidate any of their sessions (set is_active=false).
   - Forcibly-revoked sessions should survive restarts and be reflected in authentication checks (your middleware should check DB session active status for the token when present).

7) Auditing
   - Always record who did admin actions: who promoted/demoted, who revoked session X, timestamp, reason.
   - Store audit events in a separate append-only table to facilitate compliance and forensics.

8) Tests we strongly recommend
   - Unit tests for promote/demote handlers: role change + sessions invalidated + audit record created.
   - Test that duplicate login is blocked while session active, and allowed once session deactivated.
   - Admin session listing tests: only admins can access, returns correct shape, and does not leak full tokens in production mode.

9) Deploy / operational considerations
   - Add DB indexes to sessions for fast lookups by user-id and by token hash
   - Add a background job to cleanup expired/old sessions and soft-delete audit logs per retention policy

Sample handler pseudo-code (promotion + invalidate sessions)
```go
func (h *AdminHandler) PromoteUser(c *gin.Context) {
  id := parseParam(c, "id")
  // 1) update role in DB (transaction)
  if err := h.DB.Model(&entities.User{}).Where("id = ?", id).Update("role", "admin").Error; err != nil { ... }
  // 2) invalidate sessions for that user
  _ = h.DB.Model(&entities.Session{}).Where("user_id = ?", id).Update("is_active", false).Error
  // 3) write audit event: who, when, what
}
```

Next steps we can implement for you (pick any):
- Add secure admin endpoints (promote/demote + revoke session) and tests
- Replace stored tokens with token-hash storage and change middleware to look up token by hash
- Add audit table and store admin actions
- Add UI-safe token previews (mask tokens)

---

If you want, I can implement one of the next steps now — tell me which you prefer and I’ll start. If you're ready to move on to something else, tell me and I’ll wait for the next instruction.
