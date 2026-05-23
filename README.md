# Backend Assignment

## Parts

- [Part 1 — Rate-Limited API](part_1/README.md)
- [Part 2 — Product Catalog with Media](part_2/README.md)

---

## AI tools

Used **Antigravity (Claude Sonnet)** as an AI pair-programming assistant throughout this assignment.

What it helped with:
- Scaffolding the initial file structure and boilerplate
- Catching edge cases (e.g. whitespace-only `user_id`, `url.Parse` vs `url.ParseRequestURI`)
- Writing curl examples in TESTING.md
- Drafting README sections

All code was reviewed, understood, and iterated on by me. Design decisions (rolling window vs fixed, 409 for duplicate SKU, `RWMutex` for the product store, `copyProduct` for race safety, etc.) were made and reasoned through during the session.