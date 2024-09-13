# go-http

## Usage
Get the package with
```
go get github.com/karaMuha/go-http
```

Now import the package, initialize a new router and register a route with a target and handler. A target consists of a http method and a path like "POST /users". A handler is a func that accepts a `net.Conn` variable and a pointer to a `HttpRequest` variable as input 
```
package main

import (
	"fmt"
	"net"

	goHttp "github.com/karaMuha/go-http"
)

func main() {
	router := goHttp.NewRouter()

	router.HandleFunc("GET /", func(conn net.Conn, r *goHttp.Request) {
		fmt.Printf("Method: %s, Target: %s, Version: %s, Body: %s\n", r.Method, r.URI, r.HttpVersion, string(r.Body))
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	})
}
```

After that initialize a new server with a port and the newly created router and listen for incoming requests
```
package main

import (
	"fmt"
	"net"

	goHttp "github.com/karaMuha/go-http"
)

func main() {
	router := goHttp.NewRouter()

	router.HandleFunc("GET /", func(conn net.Conn, r *goHttp.Request) {
		fmt.Printf("Method: %s, Target: %s, Version: %s, Body: %s\n", r.Method, r.URI, r.HttpVersion, string(r.Body))
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	})

	server := goHttp.NewServer("8080", router)
	err := server.Listen()
	if err != nil {
		fmt.Println(err)
	}
}
```

## ToDos
- provide a developer firendly response writer instead of the raw net.Conn connection
- provide useful functionalities for HttpRequest struct
- allow path values
- allow query params
- more to come...