package header

// Type represents HTTP header fields as defined by IANA.
// Reference: https://www.iana.org/assignments/http-fields/http-fields.xhtml
type Type string

// noinspection GoUnusedConst
const (
	Accept                        Type = "Accept"
	AcceptCharset                 Type = "Accept-Charset"
	AcceptDatetime                Type = "Accept-Datetime"
	AcceptEncoding                Type = "Accept-Encoding"
	AcceptLanguage                Type = "Accept-Language"
	AcceptPatch                   Type = "Accept-Patch"
	AcceptRanges                  Type = "Accept-Ranges"
	AccessControlAllowCredentials Type = "Access-Control-Allow-Credentials"
	AccessControlAllowHeaders     Type = "Access-Control-Allow-Headers"
	AccessControlAllowMethods     Type = "Access-Control-Allow-Methods"
	AccessControlAllowOrigin      Type = "Access-Control-Allow-Origin"
	AccessControlExposeHeaders    Type = "Access-Control-Expose-Headers"
	AccessControlMaxAge           Type = "Access-Control-Max-Age"
	AccessControlRequestHeaders   Type = "Access-Control-Request-Headers"
	AccessControlRequestMethod    Type = "Access-Control-Request-Method"
	Allow                         Type = "Allow"
	Authorization                 Type = "Authorization"
	CacheControl                  Type = "Cache-Control"
	ContentDisposition            Type = "Content-Disposition"
	ContentEncoding               Type = "Content-Encoding"
	ContentLanguage               Type = "Content-Language"
	ContentLength                 Type = "Content-Length"
	ContentLocation               Type = "Content-Location"
	ContentMD5                    Type = "Content-MD5"
	ContentRange                  Type = "Content-Range"
	ContentType                   Type = "Content-Type"
	Cookie                        Type = "Cookie"
	DoNotTrack                    Type = "DNT"
	ETag                          Type = "ETag"
	Expires                       Type = "Expires"
	IfMatch                       Type = "If-Match"
	IfModifiedSince               Type = "If-Modified-Since"
	IfNoneMatch                   Type = "If-None-Match"
	IfRange                       Type = "If-Range"
	IfUnmodifiedSince             Type = "If-Unmodified-Since"
	LastModified                  Type = "Last-Modified"
	Link                          Type = "Link"
	Location                      Type = "Location"
	MaxForwards                   Type = "Max-Forwards"
	Origin                        Type = "Origin"
	P3P                           Type = "P3P"
	Pragma                        Type = "Pragma"
	ProxyAuthenticate             Type = "Proxy-Authenticate"
	ProxyAuthorization            Type = "Proxy-Authorization"
	Range                         Type = "Range"
	Referer                       Type = "Referer"
	Refresh                       Type = "Refresh"
	RetryAfter                    Type = "Retry-After"
	Server                        Type = "Server"
	SetCookie                     Type = "Set-Cookie"
	StrictTransportSecurity       Type = "Strict-Transport-Security"
	TE                            Type = "TE"
	TransferEncoding              Type = "Transfer-Encoding"
	Upgrade                       Type = "Upgrade"
	UserAgent                     Type = "User-Agent"
	Vary                          Type = "Vary"
	Via                           Type = "Via"
	WWWAuthenticate               Type = "WWW-Authenticate"
	Warning                       Type = "Warning"
	XRequestedWith                Type = "X-Requested-With"
)

// String returns the string representation of the header field.
func (t Type) String() string {
	return string(t)
}

// Parse parses the header field from a string.
func Parse(value string) Type {
	return Type(value)
}
