package go_req_res

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type httpHandler = func(w http.ResponseWriter, r *http.Request)

// Wraps an http handler and captures the response http status and body, then writes it to the given io.Writer
func CaptureResponse(dst io.Writer, h httpHandler) httpHandler {
	f := func(w http.ResponseWriter, r *http.Request) {

		// create a buffer that intercepts what the handler will write
		buffer := buffer{
			ResponseWriter: w,
		}

		// execute the handler, passing in our interceptor
		h(&buffer, r)

		// write the intercepted body and http status to the given io.Writer
		dst.Write(buffer.Bytes())
	}

	return f
}

// Wraps an http handler and captures the request http status and body, then writes it to the given io.Writer
func CaptureRequest(dst io.Writer, h httpHandler) httpHandler {
	f := func(w http.ResponseWriter, r *http.Request) {

		// read the request body's contents and print to the give io.Writer
		body, e := ioutil.ReadAll(r.Body)
		if e != nil {
			log.Println("failed to read request body, but let the actual handler still do its work")
		}
		out := fmt.Sprintf(`{"Method":"%s","Body":%s}`, r.Method, body)
		dst.Write([]byte(out))

		// set the original HttpRequest to read the body again (because it was 'consumed' already by our initial read)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// execute the handler
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
