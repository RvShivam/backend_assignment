# Part 1 — Rate-Limited API

## Running

```bash
go run .
```

Starts on port 8080. See [TESTING.md](TESTING.md) for curl examples.

---

## POST /request

```json
{
  "user_id": "alice",
  "payload": { "any": "json value" }
}
```

Both fields are required. `user_id` must be a non-empty, non-whitespace string. `payload` can be any valid JSON type (object, array, string, number, boolean).

- `201` — request accepted
- `400` — missing/empty/whitespace `user_id`, missing `payload`, or malformed JSON
- `429` — rate limit exceeded (more than 5 requests in the last 60 seconds for that user)

---

## GET /stats

Returns per-user counts:

```json
{
  "alice": {
    "accepted_in_window": 3,
    "rejected_total": 1
  }
}
```

- `accepted_in_window` — accepted requests still inside the current rolling 60-second window
- `rejected_total` — cumulative rejections since the process started (not reset per window)

---

## Design notes

**Status code** — `201` for every accepted request. The assignment allows 200 or 201; I picked 201 since each call records a new entry in the store.

**Rate limiting** — rolling window. Each accepted request stores a timestamp; anything older than 60 seconds is dropped before the count is checked. This avoids the burst problem you get at fixed-window boundaries.

**Concurrency** — a single `sync.Mutex` wraps all reads and writes. Parallel requests for the same `user_id` can't race past the limit.

---

## Production limitations

- **No shared state** — multiple instances each have their own counter, so the global rate isn't enforced across replicas. Redis (sorted set or `INCR`+`EXPIRE`) would fix this.
- **In-memory only** — everything resets on restart.
- **Memory** — the user map grows indefinitely; stale entries are never evicted.
- **No auth** — any caller can pass any `user_id`. In production the key should come from a verified token.
