package method

// Type represents HTTP methods as defined by IANA.
// Reference: https://www.iana.org/assignments/http-methods/http-methods.xhtml
type Type string

const (
	GET     Type = "GET"
	HEAD    Type = "HEAD"
	POST    Type = "POST"
	PUT     Type = "PUT"
	PATCH   Type = "PATCH"
	DELETE  Type = "DELETE"
	CONNECT Type = "CONNECT"
	OPTIONS Type = "OPTIONS"
	TRACE   Type = "TRACE"
)

// String returns the string representation of the HTTP method.
func (t Type) String() string {
	return string(t)
}

// Parse parses the string into a Type.
func Parse(value string) Type {
	return Type(value)
}
