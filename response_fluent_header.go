package fastshot

import "net/http"

type ResponseFluentHeader struct {
	header http.Header
}

func (r *Response) Header() *ResponseFluentHeader {
	return r.header
}

func (h *ResponseFluentHeader) Get(key string) string {
	return h.header.Get(key)
}

func (h *ResponseFluentHeader) GetAll(key string) []string {
	return h.header[key]
}

func (h *ResponseFluentHeader) Keys() []string {
	keys := make([]string, 0, len(h.header))
	for k := range h.header {
		keys = append(keys, k)
	}
	return keys
}
