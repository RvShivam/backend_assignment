# Testing — Part 2

Start the server:

```bash
cd part_2
go run .
```

Runs on port 8081.

---

## POST /products

**Create a product with media**

```bash
curl -s -X POST http://localhost:8081/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Widget A",
    "sku": "SKU-001",
    "image_urls": [
      "https://cdn.example.com/products/sku-001/img-1.jpg",
      "https://cdn.example.com/products/sku-001/img-2.jpg"
    ],
    "video_urls": [
      "https://cdn.example.com/products/sku-001/demo.mp4"
    ]
  }'
```

Expected: `201` with a full product body including `id`.

**Create without media (image_urls/video_urls optional)**

```bash
curl -s -X POST http://localhost:8081/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Widget B", "sku": "SKU-002"}'
```

**Duplicate sku → 409**

```bash
curl -s -X POST http://localhost:8081/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Other", "sku": "SKU-001"}'
```

**Missing name → 400**

```bash
curl -s -X POST http://localhost:8081/products \
  -H "Content-Type: application/json" \
  -d '{"sku": "SKU-003"}'
```

**Invalid URL → 400**

```bash
curl -s -X POST http://localhost:8081/products \
  -H "Content-Type: application/json" \
  -d '{"name": "X", "sku": "SKU-004", "image_urls": ["not-a-url"]}'
```

---

## GET /products

```bash
curl -s http://localhost:8081/products
```

With pagination:

```bash
curl -s "http://localhost:8081/products?limit=5&offset=0"
```

Notice the response contains `image_count` and `video_count` but no URL arrays.

---

## GET /products/:id

Copy an `id` from the create response, then:

```bash
curl -s http://localhost:8081/products/<id>
```

The response includes the full `image_urls` and `video_urls` arrays.

**Unknown id → 404**

```bash
curl -s http://localhost:8081/products/doesnotexist
```

---

## POST /products/:id/media

```bash
curl -s -X POST http://localhost:8081/products/<id>/media \
  -H "Content-Type: application/json" \
  -d '{
    "image_urls": ["https://cdn.example.com/products/sku-001/img-3.jpg"]
  }'
```

Expected: `200` with the updated product showing all URLs including the new one.

**Empty body → 400**

```bash
curl -s -X POST http://localhost:8081/products/<id>/media \
  -H "Content-Type: application/json" \
  -d '{"image_urls": [], "video_urls": []}'
```

---

## Bulk seed (optional)

To create 100 products quickly and test pagination:

**Bash**

```bash
for i in $(seq 1 100); do
  curl -s -o /dev/null -X POST http://localhost:8081/products \
    -H "Content-Type: application/json" \
    -d "{\"name\": \"Product $i\", \"sku\": \"SKU-$(printf '%04d' $i)\", \"image_urls\": [\"https://cdn.example.com/p/$i/1.jpg\", \"https://cdn.example.com/p/$i/2.jpg\"]}"
done
```

**PowerShell**

```powershell
1..100 | ForEach-Object {
  $body = @{
    name = "Product $_"
    sku = "SKU-$($_.ToString('0000'))"
    image_urls = @("https://cdn.example.com/p/$_/1.jpg", "https://cdn.example.com/p/$_/2.jpg")
  } | ConvertTo-Json
  Invoke-RestMethod -Uri http://localhost:8081/products -Method POST -ContentType "application/json" -Body $body | Out-Null
}
```

Then verify pagination returns only summaries (no URL arrays):

```bash
curl -s "http://localhost:8081/products?limit=10&offset=0"
```
