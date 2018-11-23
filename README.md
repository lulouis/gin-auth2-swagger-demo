# Important Step one by one

Gen doc

```console
$ go get -u github.com/swaggo/swag/cmd/swag
$ swag init
```

Run app

```console
$ go run main.go
```

Build
```console
set CGO_ENABLED=0 GOOS=linux 
go build -a -installsuffix cgo -o gin-auth2-swagger-demo .
```

Build image in Mac
```console
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gin-auth2-swagger-demo .
```


Run
```console
 docker run -itd -p 8080:8080 --rm --name demo gin-auth2-swagger-demo
```

[open swagger](http://localhost:8080/swagger/index.html)

