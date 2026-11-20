# Prompts

I used Claude Code (Claude Opus 5.5) in the terminal. I had only used Go for a few small things about two years ago, so I mostly used it as a teacher and reviewer: it refreshed the Go concepts I needed, I wrote the code, then it reviewed my code. For some parts I asked it to write code or tests, and those are marked below.

These are the most important prompts, slightly reworded so they're easier to read. Short questions about Go syntax are left out.

## Summary

| Part | Written by |
|---|---|
| `calc` package (operations, registry, `Calculate`) | me |
| `calc` tests | me, extra cases added by AI |
| `api` handler | me, reviewed by AI |
| `api` tests | AI |
| `cmd/server/main.go` | AI |

## 1. Plan

**Prompt:** (pasted the assignment) "This is the project. We need a plan first. I haven't used Go much in recent years, so explain the Go parts as we go and let me write the core code myself."

**Result:** A plan: Go standard library `net/http`, a pure `calc` package separate from the HTTP layer, one `POST /api/v1/calculate` endpoint, a single error format, table-driven tests, Vite + React + TypeScript, and a multi-stage Dockerfile.

**What I did:** Kept the overall structure. Chose float64 instead of a decimal library to keep the scope small (see README).

## 2. Review of the first operations

**Prompt:** "I wrote Add, Subtract, Multiply and Divide with a table-driven test. Review my files like a senior Go reviewer would."

**Result:** It pointed out a typo in an exported error name, a `switch` inside the test (logic in a test), and comparing errors with `!=`.

**What I did:** Fixed the typo, changed the test table to store the function itself (`fn func(a, b float64) float64`), and switched to `errors.Is`.

## 3. Parentheses

**Prompt:** "Should we support parentheses?"

**Result:** It suggested not to. Expressions need a parser and a different API, and the brief says to prioritize correctness over extra features.

**What I did:** Left it out and listed it in the README as a possible next step.

## 4. Tests for the calc package

**Prompt:** "I have implemented Power, Square Root and Percent, and a Calculate function that looks up an operation by name. Please add test cases for the relevant edge cases."

**Result:** Added cases for each operation, unknown operations, wrong operand counts, division by zero and NaN/Inf results. It noted that `SquareRoot(-4)` returned NaN without an error. 4 of the new tests failed.

**What I did:** Changed `SquareRoot` to return `ErrNegativeSqrt`. The failing tests showed a bug in my NaN/Inf check (`err != nil` instead of `err == nil`), so the check never ran. I fixed it.

## 5. Server

**Prompt:** "Write main.go for me. It should read the port from an environment variable and start the server with the router from the api package."

**Result:** `cmd/server/main.go` with a `PORT` env variable, server timeouts and graceful shutdown.

**What I did:** Reviewed it and went through the explanation of the timeouts, goroutine and context parts.

## 6. Review of the API

**Prompt:** "Check what I've done so far in the api package. Run the server and try the endpoints."

**Result:** It tested the server with curl and found that unknown operation, wrong operand count and negative square root returned 500 instead of a client error, and that the health route used POST.

**What I did:** Mapped every `calc` error to a status code. I chose 422 for valid JSON that can't be calculated and 400 for JSON that can't be read. Wrapped errors in `{"error": {...}}`.

## 7. API tests

**Prompt:** "Add tests for the handler. Cover success, invalid input and every error code."

**Result:** 26 handler tests (success, invalid JSON, body over 1 MB, every error code, wrong HTTP method, health check). One failed because my health route was still `POST`.

**What I did:** Changed it to `GET /healthz`.

## 8. Scaffolding frontend
**Prompt:**: Scaffold the frontend with Vite, React, TypeScript and Vitest, and add a dev proxy to the Go API

