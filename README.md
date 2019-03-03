# go-req-res

Output the request and response body of an http handler
to an io.Writer of your choice.

So you keep your http handler clean and focused
on the task at hand.

## Usage

Capturing request and printing to stdout

```cgo
go_req_res.CaptureRequest(os.Stdout, my_handler_func)
```

Capturing response and printing to stdout

```cgo
go_req_res.CaptureResponse(os.Stdout, my_handler_func)
```

Capturing both request and response and printing to stdout

```cgo
go_req_res.CaptureResponse(os.Stdout, CaptureRequest(os.Stdout, my_handler_func))
```

## Unit Tests

Just run `go test -cover`

## Author

Jonathan Lopez

https://github.com/jonathangersam