# Part 2 — Product Catalog with Media

## Running

```bash
cd part_2
go run .
```

Starts on port 8081. See [TESTING.md](TESTING.md) for curl examples.

---

## Endpoints

### POST /products

Creates a product.

```json
{
  "name": "Widget A",
  "sku": "SKU-001",
  "image_urls": ["https://cdn.example.com/products/sku-001/img-1.jpg"],
  "video_urls": ["https://cdn.example.com/products/sku-001/demo.mp4"]
}
```

`name` and `sku` are required. `image_urls` and `video_urls` are optional.

- `201` — created; body is the full product including the assigned `id`
- `400` — missing/empty name or sku, invalid URLs, too many URLs
- `409` — sku already exists

I chose 409 over 400 for duplicate sku because it's a conflict with existing state, not a malformed request.

---

### GET /products

Lists products. Returns summaries only — no URL arrays.

```
GET /products?limit=20&offset=0
```

Default limit: 20. Max limit: 100. Default offset: 0.

```json
{
  "data": [
    {
      "id": "a1b2c3d4e5f6a7b8",
      "name": "Widget A",
      "sku": "SKU-001",
      "image_count": 2,
      "video_count": 1,
      "thumbnail_url": "https://cdn.example.com/products/sku-001/img-1.jpg",
      "created_at": "2026-05-22T10:00:00Z"
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

`thumbnail_url` is the first image URL for that product, omitted if none.

---

### GET /products/:id

Returns the full product including all `image_urls` and `video_urls`.

- `404` — unknown id

---

### POST /products/:id/media

Appends URLs to an existing product.

```json
{
  "image_urls": ["https://cdn.example.com/products/sku-001/img-3.jpg"],
  "video_urls": []
}
```

At least one of `image_urls` or `video_urls` must be provided and non-empty.

- `200` — body is the updated full product
- `400` — no URLs provided, or validation failure
- `404` — unknown id

---

## Validation rules

- `name` and `sku`: required, non-empty after trimming whitespace
- URLs: must start with `http://` or `https://`, have a non-empty host, and be at most 2048 characters
- Max 20 URLs per array per request (applies to both `image_urls` and `video_urls` independently)

---

## Data model and performance

Products are stored as `map[string]*Product` (id → product) with a separate `map[string]string` for sku → id lookups and a `[]string` slice to maintain insertion order for pagination. Media URLs are embedded directly in the product struct.

**List vs detail:** `GET /products` builds lightweight `ProductSummary` DTOs that only read `len(p.ImageURLs)` and `len(p.VideoURLs)` — the actual URL strings are never iterated or serialized. With 1,000 products each having 10 images, `GET /products?limit=20` touches 20 products and serializes zero URL strings. `GET /products/:id` serializes all URLs for that one product only.

**With PostgreSQL + CDN in production:**
- Products table: `(id, name, sku, created_at)` — indexed on sku
- Media table: `(id, product_id, type, url, position)` — indexed on product_id
- List query joins products only, no media rows fetched; image_count/video_count come from a precomputed column or a fast COUNT subquery
- Detail query fetches the product row and all its media rows
- Thumbnail is either stored on the product row directly, or derived from the first media row with `ORDER BY position LIMIT 1`
- CDN URLs are stored as strings; the API never generates or validates them against the CDN

---

## Production limitations

- **In-memory only** — state is lost on restart
- **No pagination on media** — a product with thousands of URLs returns them all in the detail response
- **No auth or ownership** — any caller can modify any product
- **Single process** — the in-memory store can't be shared across replicas
