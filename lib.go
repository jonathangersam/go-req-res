package go_req_res

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

type httpHandler = func(w http.ResponseWriter, r *http.Request)

// Wraps an http handler and captures the response http status and body, then writes it to the given io.Writer
func CaptureResponse(dst io.Writer, h httpHandler) httpHandler {
	f := func(w http.ResponseWriter, r *http.Request) {
		buffer := buffer{
			ResponseWriter: w,
		}

		h(&buffer, r)

		dst.Write(buffer.Bytes())
	}

	return f
}

// ...
func CaptureRequest(dst io.Writer, h httpHandler) httpHandler {
	f := func(w http.ResponseWriter, r *http.Request) {

		dump, _ := httputil.DumpRequest(r, true)
		dst.Write(dump)
		//fmt.Fprintf(dst, "%q", dump)

		h(w, r)

	}

	return f
}

type buffer struct {
	http.ResponseWriter

	capturedBody       []byte
	capturedStatusCode int
}

func (b *buffer) Write(p []byte) (int, error) {
	for _, el := range p {
		b.capturedBody = append(b.capturedBody, el)
	}

	return b.ResponseWriter.Write(p)
}

func (b *buffer) WriteHeader(statusCode int) {
	b.capturedStatusCode = statusCode

	b.ResponseWriter.WriteHeader(statusCode)
}

func (b *buffer) Bytes() []byte {
	s := fmt.Sprintf(`{"HttpStatus":%d,"Body":%s}`, b.capturedStatusCode, b.capturedBody)
	return []byte(s)
}
