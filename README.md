<div style="text-align: center;">
    <h1>Fast Shot</h1>
    <p>A Fluent Go REST Client Library</p>
    <p>
        <a href="https://goreportcard.com/report/opus-domini/fast-shot"><img src="https://goreportcard.com/badge/opus-domini/fast-shot" alt="Go Report Badge"></a>
        <a href="https://godoc.org/github.com/opus-domini/fast-shot"><img src="https://godoc.org/github.com/opus-domini/fast-shot?status.svg" alt="Go Doc Badge"></a>    
        <a href="https://github.com/opus-domini/fast-shot/actions/workflows/coverage.yml"><img src="https://github.com/opus-domini/fast-shot/actions/workflows/coverage.yml/badge.svg" alt="Converage Actions Badge"></a>
        <a href="https://codecov.io/gh/opus-domini/fast-shot"><img src="https://codecov.io/gh/opus-domini/fast-shot/graph/badge.svg?token=C80QDL5W7T" alt="Codecov Badge"/></a>        
        <a href="https://github.com/opus-domini/fast-shot/blob/main/LICENSE"><img src="https://img.shields.io/github/license/opus-domini/fast-shot.svg" alt="License Badge"></a>
        <a href="https://github.com/avelino/awesome-go"><img src="https://awesome.re/mentioned-badge.svg" alt="Mentioned in Awesome Go"></a>
    </p>
</div>

Fast Shot is a robust, feature-rich, and highly configurable HTTP client for Go. Crafted with modern Go practices in mind, it offers a fluent, chainable API that allows for clean, idiomatic code.

## Why Fast Shot?

* **Fluent & Chainable API**: Write expressive, readable, and flexible HTTP client code.
* **Ease of Use**: Reduce boilerplate, making HTTP requests as straightforward as possible.
* **Rich Features**: From headers to query parameters and JSON support, Fast Shot covers your needs.

## Features 🌟

* Fluent API for HTTP requests
* Extensible authentication
* Customizable HTTP headers
* Query parameter manipulation
* JSON request and response support
* Built-in error handling
* Well-tested

## Installation 🔧

To install Fast Shot, run the following command:

```bash 
go get github.com/opus-domini/fast-shot
```

## Quick Start 🚀

Here's how you can make a simple GET request using Fast Shot:

```go
package main

import (
    "fmt"
    fastshot "github.com/opus-domini/fast-shot"
)

func main() {
    client := fastshot.NewClient("https://api.example.com").
        Auth().
            BearerToken("your_token_here").
        End().        
        Build()

    payload := map[string]interface{}{
        "key1": "value1",
        "key2": "value2",
    }
	
    response, err := client.POST("/endpoint").
        BodyJSON(payload).
        Send()
	
    if err != nil {
        fmt.Println("Error: ", err)
        return
    }
	
    // Process response...
}
```

## Advanced Usage 🛠️

### Fluent API

Easily chain multiple settings in a single line:

```go 
client := fastshot.NewClient("https://api.example.com").
    Auth().BearerToken("your-bearer-token").End().
    Header().Add("My-Header", "My-Value").End().
    Config().SetTimeout(time.Second * 30).End().
    Build()
```

### Authentication

Fast Shot supports various types of authentication:

```go
// Bearer Token
builder.Auth().BearerToken("your-bearer-token").End()

// Basic Authentication
builder.Auth().BasicAuth("username", "password").End()
```

### Custom Headers and Cookies

Add your own headers and cookies effortlessly:

```go 
// Add Custom Headers
builder.
    Header().
        Add("Custom-Header-1", "value").
        Add("Custom-Header-2", "value").
    End()

// Add Custom Cookie
builder.
    Cookie().
        Add(&http.Cookie{Name: "session_id", Value: "id"}).
    End()
```

### Advanced Configurations

Control every aspect of the HTTP client:

```go
// Set Timeout
builder.Config().SetTimeout(time.Second * 30).End()

// Custom Transport
builder.Config().SetCustomTransport(myCustomTransport).End()
````

## Contributing 🤝

Your contributions are always welcome! Feel free to create pull requests, submit issues, or contribute in any other way.
