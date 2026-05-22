# Testing — Part 1

Start the server first:

```bash
cd part_1
go run .
```

---

## POST /request

**Happy path**

```bash
curl -s -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"user_id": "alice", "payload": {"action": "buy"}}'
```

Expected: `201`
```json
{"message": "request accepted"}
```

---

**Missing user_id**

```bash
curl -s -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"payload": {"action": "buy"}}'
```

Expected: `400`
```json
{"error": "user_id cannot be empty"}
```

---

**Whitespace-only user_id**

```bash
curl -s -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"user_id": "   ", "payload": 1}'
```

Expected: `400`
```json
{"error": "user_id cannot be empty"}
```

---

**Missing payload**

```bash
curl -s -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"user_id": "alice"}'
```

Expected: `400`
```json
{"error": "user_id and payload are required"}
```

---

**Malformed JSON**

```bash
curl -s -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d 'not json'
```

Expected: `400`
```json
{"error": "user_id and payload are required"}
```

---

## Rate limiting

Send 6 requests for the same user in quick succession. The first 5 should be accepted, the 6th rejected.

**Bash (Linux/macOS)**

```bash
for i in $(seq 1 6); do
  curl -s -o /dev/null -w "%{http_code}\n" -X POST http://localhost:8080/request \
    -H "Content-Type: application/json" \
    -d '{"user_id": "bob", "payload": true}'
done
```

Expected output:
```
201
201
201
201
201
429
```

**PowerShell (Windows)**

```powershell
1..6 | ForEach-Object {
  $r = Invoke-WebRequest -Uri http://localhost:8080/request -Method POST `
    -ContentType "application/json" `
    -Body '{"user_id": "bob", "payload": true}' `
    -ErrorAction SilentlyContinue
  $r.StatusCode
}
```

After waiting 60 seconds, the window resets and the next request will be accepted again:

```bash
sleep 61
curl -s -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"user_id": "bob", "payload": "retry"}'
# → 201
```

---

## GET /stats

```bash
curl -s http://localhost:8080/stats
```

Example response after the rate-limit test above:

```json
{
  "alice": {
    "accepted_in_window": 1,
    "rejected_total": 0
  },
  "bob": {
    "accepted_in_window": 5,
    "rejected_total": 1
  }
}
```

`accepted_in_window` drops to 0 once the 60-second window passes. `rejected_total` is cumulative and does not reset.
