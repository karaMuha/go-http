# go-http

## Usage
Get the package with
```
go get github.com/karaMuha/go-http
```

Now import the package, initialize a new router and register a route with a target and handler. A target consists of a http method and a path like "POST /users". A handler is a func that accepts a pointer to a `HttpResponse` variable and a pointer to a `HttpRequest` variable as input 
```
package main

import (
	"fmt"
	"net"

	goHttp "github.com/karaMuha/go-http"
)

func main() {
	router := goHttp.NewRouter()

	router.HandleFunc("GET /", func(res *goHttp.HttpResponse, req *goHttp.HttpRequest) {
		fmt.Printf("Method: %s, Target: %s, Version: %s, Body: %s\n", req.Method, req.URI, req.HttpVersion, string(req.Body))
		
		cookie := &goHttp.Cookie{
			Name: "TestKey",
			Value: "TestValue"
			Secure: true,
			HttpOnly: true,
			Expires: time.Now().Add(1 * time.Hour)
		}

		res.SetCookie(cookie)
		res.SetHeader("Content-Type", "application/json")
		res.WriteStatus(200)
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

	router.HandleFunc("GET /", func(res *goHttp.HttpResponse, req *goHttp.HttpRequest) {
		fmt.Printf("Method: %s, Target: %s, Version: %s, Body: %s\n", req.Method, req.URI, req.HttpVersion, string(req.Body))
		
		cookie := &goHttp.Cookie{
			Name: "TestKey",
			Value: "TestValue"
			Secure: true,
			HttpOnly: true,
			Expires: time.Now().Add(1 * time.Hour)
		}

		res.SetCookie(cookie)
		res.SetHeader("Content-Type", "application/json")
		res.WriteStatus(200)
	})

	server := goHttp.NewServer("8080", router)
	err := server.Listen()
	if err != nil {
		fmt.Println(err)
	}
}
```

## ToDos
- provide useful functionalities for HttpResponse struct
- provide useful functionalities for HttpRequest struct
- parse cookies from request
- allow path values
- allow query params
- more to come...