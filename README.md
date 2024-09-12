# go-http

## Usage
Get the package with
```
go get github.com/karaMuha/go-http
```

Now import the package, initialize a new router and register a route
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

After that initialize a new server with a port and the newly created router
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