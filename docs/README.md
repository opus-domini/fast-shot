<div align="center">
    <img src="assets/images/logo.png" alt="Logo fast-shot" width="320"/>
    <hr />
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

## Table of Contents

- [Why Fast Shot?](#why-fast-shot)
- [Features](#features-)
- [Installation](#installation-)
- [Quick Start](#quick-start-)
- [Advanced Usage](#advanced-usage-)
- [Contributing](#contributing-)

## Why Fast Shot?

* **Fluent & Chainable API**: Write expressive, readable, and flexible HTTP client code.
* **Ease of Use**: Reduce boilerplate, making HTTP requests as straightforward as possible.
* **Rich Features**: From headers to query parameters and JSON support, Fast Shot covers your needs.
* **Advanced Retry Mechanism**: Built-in support for retries with various backoff strategies.

## Features 🌟

* Fluent and chainable API for clean, expressive code
* Comprehensive HTTP method support (GET, POST, PUT, DELETE, etc.)
* Flexible authentication options (Bearer Token, Basic Auth, Custom)
* Easy manipulation of headers, cookies, and query parameters
* Advanced retry mechanism with customizable backoff strategies
* Client-side load balancing for improved reliability
* JSON request and response support
* Timeout and redirect control
* Proxy support
* Extensible and customizable for specific needs
* Well-tested and production-ready

## Installation 🔧

To install Fast Shot, run the following command:

```bash 
go get github.com/opus-domini/fast-shot
```

## Quick Start 🚀

Here's how you can make a simple POST using Fast Shot:

```go
package main

import (
    "fmt"
    fastshot "github.com/opus-domini/fast-shot"
    "github.com/opus-domini/fast-shot/constant/mime"
)

func main() {
    client := fastshot.NewClient("https://api.example.com").
        Auth().BearerToken("your_token_here").
        Build()

    payload := map[string]interface{}{
        "key1": "value1",
        "key2": "value2",
    }

    response, err := client.POST("/endpoint").
        Header().AddAccept(mime.JSON).
        Body().AsJSON(payload).
        Send()

    // Check for request send problems.
    if err != nil {
        panic(err) // (¬_¬")
    }

    // Check for (4xx || 5xx) errors response.
    if response.Status().IsError() {
        panic(response.Body().AsString()) // ¯\_(ツ)_/¯
    }
	
    var result map[string]interface{}
    _ := response.Body().AsJSON(&result)

    // Congrats! Do something awesome with the result (¬‿¬)
}
```

## Advanced Usage 🤖

### Fluent API

Easily chain multiple settings in a single line:

```go 
client := fastshot.NewClient("https://api.example.com").
    Auth().BearerToken("your-bearer-token").
    Header().Add("My-Header", "My-Value").
    Config().SetTimeout(30 * time.Second).
    Build()
```
### Advanced Retry Mechanism

Handle transient failures with customizable backoff strategies:

```go 
client.POST("/resource").
    Retry().SetExponentialBackoff(2 * time.Second, 5, 2.0).
    Send()
```

This new retry feature supports:
- Constant backoff
- Exponential backoff
- Full jitter for both constant and exponential backoff
- Custom retry conditions
- Maximum delay setting

### Out-of-the-Box Support for Client Load Balancing

Effortlessly manage multiple endpoints:

```go
client := fastshot.NewClientLoadBalancer([]string{
    "https://api1.example.com",
    "https://api2.example.com",
    "https://api3.example.com",
    }).
    Config().SetTimeout(time.Second * 10).
    Build()
```

This feature allows you to distribute network traffic across several servers, enhancing the performance and reliability of your applications.

### Authentication

Fast Shot supports various types of authentication:

```go
// Bearer Token
builder.Auth().BearerToken("your-bearer-token")

// Basic Authentication
builder.Auth().BasicAuth("username", "password")

// Custom Authentication
builder.Auth().Set("custom-authentication-header")
```

### Custom Headers and Cookies

Add your own headers and cookies effortlessly:

```go 
// Add Custom Header
builder.Header().
    Add("header", "value")

// Add Multiple Custom Headers
builder.Header().
    AddAll(map[string]string{
        "key1": "value1",
        "key2": "value2",
        "key3": "value3",
    })

// Add Custom Cookie
builder.Cookie().
    Add(&http.Cookie{Name: "session_id", Value: "id"})
```

### Advanced Configurations

Control every aspect of the HTTP client:

```go
// Set Timeout
builder.Config().
    SetTimeout(time.Second * 30)

// Set Follow Redirect
builder.Config().
    SetFollowRedirects(false)

// Set Custom Transport
builder.Config().
    SetCustomTransport(myCustomTransport)

// Set Proxy
builder.Config().
    SetProxy("http://my-proxy-server:port")
```

### Response Handling

Extract information from the response with ease:

```go
// Fluent response status check
response.Status()

// Easy response body access and conversion
response.Body()

// Get response headers to inspect
response.Header()

// Get response cookies for further processing 
response.Cookie()

// Get raw response if needed
response.Raw()

// and more...
```

## Contributing 🤝

We welcome contributions to Fast Shot! Here's how you can contribute:

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Write your code and tests
4. Ensure all tests pass
5. Submit a pull request

Please make sure to update tests as appropriate and adhere to the existing coding style.

For more detailed information, check out our [CONTRIBUTING.md](https://github.com/opus-domini/fast-shot/blob/main/CONTRIBUTING.md) file.

## Stargazers over time ⭐

[![Stargazers over time](https://starchart.cc/opus-domini/fast-shot.svg?variant=adaptive)](https://starchart.cc/opus-domini/fast-shot)
