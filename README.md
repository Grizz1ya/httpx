# httpx

A simple, chainable HTTP client for Go with session support.

## Installation

```bash
go get github.com/Grizz1ya/httpx
```

## Features

- Session management with cookie persistence
- Chainable request building
- Support for different request formats (JSON, form data, raw body)
- Proxy support
- Static and request-specific headers

## Usage

### Basic Requests

Make a simple GET request:

```go
import "github.com/Grizz1ya/httpx"

// Simple GET request
response, err := httpx.Get("https://example.com").Do()
if err != nil {
    // handle error
}
fmt.Println(response.Text())

// GET with query parameters
response, err = httpx.Get("https://example.com").Params(map[string]interface{}{
    "page": "1",
    "limit": "10",
}).Do()

// POST with JSON body
response, err = httpx.Post("https://example.com/api").Json(map[string]interface{}{
    "username": "user",
    "password": "pass",
}).Do()
```

### Session Management

Create a session to reuse cookies and headers across requests:

```go
// Create a new session
session := httpx.NewSession()

// Add static headers that will be used in all requests
session.AddStaticHeader("User-Agent", "My Custom User Agent")

// Add cookies
session.AddCookie("example.com", "session_id", "abc123")

// Make a request using the session
response, err := session.Get("https://example.com").Do()

// Remove headers or cookies when needed
session.RemoveStaticHeader("User-Agent")
session.RemoveCookie("example.com", "session_id")
```

### Using Proxies

```go
// Create a proxy from URL string
proxy, err := httpx.NewProxyFromLine("http://127.0.0.1:8080")
if err != nil {
    // handle error
}

// Set proxy to a session
session := httpx.NewSession()
err = session.SetProxy(proxy)
if err != nil {
    // handle error
}

// Remove proxy
session.SetProxy(nil)
```

### Request Customization

```go
// Add headers to a specific request
response, err := session.Get("https://example.com").Headers(map[string]string{
    "X-Custom-Header": "value",
}).Do()

// Add form data to a POST request
response, err = session.Post("https://example.com/form").Data(map[string]interface{}{
    "field1": "value1",
    "field2": "value2",
}).Do()

// Add JSON to a POST request
response, err = session.Post("https://example.com/api").Json(map[string]interface{}{
    "key": "value",
}).Do()

// Set a raw body
response, err = session.Post("https://example.com/api").Body([]byte("raw body content")).Do()
```

### Response Handling

```go
// Get response as text
response, _ := httpx.Get("https://example.com").Do()
textContent := response.Text()

// Parse JSON response
jsonResponse := response.Json()
value := jsonResponse.Get("key").String()

// Access cookies
cookies := response.Cookies()
```

## Example

```go
package main

import (
    "fmt"
    "github.com/Grizz1ya/httpx"
)

func main() {
    // Create a session
    session := httpx.NewSession()
    session.AddStaticHeader("User-Agent", "Mozilla/5.0")

    // Set up a proxy
    proxy, err := httpx.NewProxyFromLine("http://127.0.0.1:8080")
    if err == nil {
        session.SetProxy(proxy)
    }

    // Make a request
    response, err := session.Get("https://api.example.com/data").Params(map[string]interface{}{
        "limit": "10",
    }).Do()

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Parse JSON response
    json := response.Json()
    items := json.GetArray("items")

    for _, item := range items {
        fmt.Println(item.Get("name").String())
    }
}
```

## License

MIT License

Copyright (c) [2025] [Grizz1ya]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
