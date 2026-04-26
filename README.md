# url-shortener

A URL shortener that fits in three files. Shorten a URL, get a slug, redirect
with HTTP 302. Slugs persist to a JSON file between restarts.

## run

```bash
go run .
```

```bash
PORT=8080 BASE_URL=https://yourdomain.com DB_PATH=urls.json go run .
```

## shorten

```bash
curl -s -X POST localhost:8080/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url": "https://github.com/AstorisTheBrave"}' | jq
# { "slug": "aB3kR9", "short": "http://localhost:8080/aB3kR9" }
```

## redirect

```bash
curl -L localhost:8080/aB3kR9
# follows the redirect to the original URL
```

## what this shows

- `net/http` handler interface (`ServeHTTP`)
- `sync.RWMutex` for safe concurrent map access
- `crypto/rand` for unguessable slugs (not `math/rand`)
- JSON file persistence with atomic read-lock during serialization
- environment-based config with sensible defaults
