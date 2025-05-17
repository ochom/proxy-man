# Proxy Man

A simple HTTP proxy server built with [Fiber](https://gofiber.io/) in Go. This app exposes a single `/proxy` endpoint that forwards HTTP requests to a target URL, acting as a programmable proxy. It is useful for debugging, testing, or as a lightweight API gateway.

## Features
- Accepts any HTTP method (GET, POST, PUT, DELETE)
- Forwards headers and body to the target URL
- Returns the proxied response
- Simple validation and error handling

## Requirements
- Go 1.18+

## Installation

Clone the repository and install dependencies:

```bash
git clone https://github.com/ochom/proxy-man.git
cd proxy-man
go mod tidy
```

## Usage

Run the server:

```bash
go run main.go
```

The server will start on the default Fiber port (3000). You can change the port in the `main.go` file if needed.

## API

### POST /proxy

Proxy an HTTP request to a target URL.

#### Request Body

```json
{
  "method": "GET|POST|PUT|DELETE",
  "url": "https://target.url/api",
  "headers": { "Header-Name": "value" },
  "body": "<raw bytes, optional>"
}
```

- `method`: HTTP method (required)
- `url`: Target URL to proxy to (required)
- `headers`: Map of headers (optional)
- `body`: Raw request body as bytes (optional)

#### Example cURL

```bash
curl -X POST http://localhost:3000/proxy \
  -H "Content-Type: application/json" \
  -d '{
    "method": "GET",
    "url": "https://httpbin.org/get",
    "headers": {"X-Test": "123"}
  }'
```

#### Response

- On success: Returns the proxied response body as JSON.
- On error: Returns `{ "error": "<error message>" }` with status 400.

## License

MIT
