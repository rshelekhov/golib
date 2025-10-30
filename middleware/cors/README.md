# CORS Middleware

HTTP middleware for handling Cross-Origin Resource Sharing (CORS) headers.

## Features

- Configurable allowed origins (supports wildcard `*`)
- Handles preflight OPTIONS requests
- Sets standard CORS headers
- Configurable allowed methods and headers

## Usage

```go
import "github.com/rshelekhov/golib/middleware/cors"

// Allow specific origins
origins := []string{"https://example.com", "https://app.example.com"}
handler := cors.Middleware(origins)(yourHandler)

// Allow all origins (development only)
handler := cors.Middleware([]string{"*"})(yourHandler)
```

## Configuration

The middleware sets the following headers:

- `Access-Control-Allow-Origin`: Set based on the origin and allowed origins list
- `Access-Control-Allow-Methods`: `GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers`: `Content-Type, Authorization`

Preflight OPTIONS requests are automatically handled with a 200 OK response.
