# Calculator App

A simple full-stack calculator built with React, TypeScript and Go as part of the Sezzle technical assessment. The frontend sends every calculation to a Go REST API.

## Run it

**With Docker** (only Docker needed):

```sh
make docker-run        # or: docker build -t sezzle-calculator . && docker run --rm -p 8080:8080 sezzle-calculator
```

Open http://localhost:8080.

**Locally** (Go 1.27+, Node 22+):

```sh
cd frontend && npm install && cd ..
make dev               # API on :8080, frontend on http://localhost:5173
```

Or in two terminals: `make run-backend` and `make run-frontend`. In development, Vite forwards `/api` to the Go server, so no CORS setup is needed.

## Tests and coverage

```sh
make test              # Go and frontend unit tests
make coverage          # reports in coverage/backend.html and coverage/frontend/index.html
make lint              # gofmt, go vet, tsc, oxlint
```

| Part | Coverage |
|---|---|
| `backend/internal/calc` | 100% |
| `backend/internal/api` | 97% |
| frontend (lines) | 93% |

## API

`POST /api/v1/calculate`

```sh
curl -X POST localhost:8080/api/v1/calculate -d '{"operation":"div","operands":[10,4]}'
# 200 {"result":2.5}

curl -X POST localhost:8080/api/v1/calculate -d '{"operation":"sqrt","operands":[9]}'
# 200 {"result":3}

curl -X POST localhost:8080/api/v1/calculate -d '{"operation":"div","operands":[10,0]}'
# 422 {"error":{"code":"DIVISION_BY_ZERO","message":"divided by zero"}}
```

| Operation | Operands | Result |
|---|---|---|
| `add`, `sub`, `mult`, `div` | 2 | a + b, a − b, a × b, a ÷ b |
| `exp` | 2 | a to the power of b |
| `sqrt` | 1 | square root of a |
| `perc` | 1 | a / 100 |

| Status | Code | When |
|---|---|---|
| 400 | `INVALID_JSON` | the body isn't valid JSON, has unknown fields, wrong types, or is over 1 MB |
| 422 | `UNKNOWN_OPERATION` | the operation name doesn't exist |
| 422 | `INVALID_OPERAND_COUNT` | wrong number of operands for the operation |
| 422 | `DIVISION_BY_ZERO` | dividing by 0 |
| 422 | `INVALID_RESULT` | square root of a negative number, or the result is too large (Inf/NaN) |

`GET /healthz` returns `{"status":"ok"}`.

## Design decisions

- **Standard library only in Go.** `net/http` routing (Go 1.22+) was enough for two routes, so no framework.
- **`calc` package has no HTTP code.** All operations are in a registry (`map[string]CalcOperation`). Adding an operation means adding one entry and its tests, without a new route. The HTTP layer only decodes JSON, calls `calc.Calculate` and maps errors to status codes.
- **400 vs 422.** 400 means the request can't be read. 422 means it can be read but can't be calculated.
- **One error format** (`{"error": {"code", "message"}}`), so the frontend can check one field and show the message.
- **float64.** Simple and enough for a calculator. Results that can't be represented (NaN, ±Inf) are returned as errors instead of invalid JSON. For money I would use a decimal type.
- **Frontend logic in one hook** (`useCalculator`), tested without rendering. Components only display state and pass key presses.
- **One container.** The Go server also serves the built frontend when `STATIC_DIR` is set, so Docker runs everything on one port.

## Assumptions

- Percent converts a number to a fraction: 50 % → 0.5.
- No parentheses or operator precedence. It works like a simple phone calculator: `2 + 3 × 4 =` gives 20.
- Operations can be chained: `2 + 3 ×` shows 5 and waits for the next number.
- Pressing two operators in a row keeps the last one. `=` without a second number does nothing.
- After a result, a digit starts a new number and an operator continues from the result.
- After an error, a digit starts fresh and operators are ignored. `C` resets everything.
- Results are rounded to 12 significant digits for display, so `0.1 + 0.2` shows `0.3`.
- Keyboard: `0–9 . + - * / ^ % r (√) Enter Esc`.

## What I'd do with more time

- Expression support (`(2+3)*4`) with a small parser that reuses the `calc` registry.
- Decimal arithmetic for exact results.
- Backspace key and calculation history.
- Request logging, rate limiting and CI.

## Use of AI

I used Claude Code. The prompts, and what I changed after each one, are in [PROMPTS.md](PROMPTS.md).
