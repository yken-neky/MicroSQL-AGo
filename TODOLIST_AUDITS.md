# TODOLIST_AUDITS

This file collects recommendations for improvements, tests and follow-ups related to the `/api/db/:manager/audits/execute` endpoint and the audit use-case implementation.

## Summary of recommendations

1. HTTP status codes
   - Use more specific status codes instead of returning 500 for all business errors. Examples:
     - 400 Bad Request — invalid request body / parse errors
     - 404 Not Found — no active connection for user+manager
     - 422 Unprocessable Entity — no scripts found for given IDs or validation errors
     - 409 Conflict — maybe for overlapping audit requests or unique constraint conflicts
     - 500 Internal Server Error — only for unexpected runtime errors

2. Connection lifecycle
   - Ensure that `sqlService.Connect()` resources are explicitly closed or pooled correctly. If `Connect` returns a resource that must be closed, call `Close(db)` in a defer or after execution.

3. Timeouts and context cancellation
   - Respect `ctx` deadlines / cancellation across query execution. If scripts are long-running, ensure the SQL execution uses context-aware clients and aborts properly when the request times out.

4. Logging and observability
   - Add structured logs around key events: connection selection, connection errors, script validation failures, script execution errors, durations
   - Consider instrumenting metrics (execution time, success/failure counters)

5. Input validation
   - Validate `database` field presence or provide a clear fallback. Sanitize / constrain accepted values.
   - Add limits for number of scripts in a single run (e.g. max 200) to avoid DoS.

6. Mode (partial/full)
   - Implement an explicit `mode` for full audit runs (e.g., `mode: full`) to run all controls. Today the code always sets `mode` to "partial".

7. Tests / E2E coverage
   - Add integration tests for: open connection → execute audit → verify created AuditRun + AuditScriptResults.
   - Add tests for negative cases: invalid SQL, missing connection, script validation failure.

8. Return structure improvements
   - Include a run status / final audit_run object in the POST response for convenience, or at minimum return `audit_run_id` (already supported when persisted).

9. Rate limiting / throttling
   - Protect endpoint from repeated heavy executions by adding rate limits for users or for same user/manager.


## Next steps (optional)
- Add an integration test suite which boots a test DB, creates required fixtures, runs full flow and asserts records persisted in `audit_runs` and `audit_script_results`.
- Implement explicit connection close to prevent resource leaks.
- Improve status codes and add tests verifying behaviour.

---

Saved recommendations on: 2025-11-27
