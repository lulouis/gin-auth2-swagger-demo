FROM alpine:latest

COPY ./gin-auth2-swagger-demo /demo/gin-auth2-swagger-demo

EXPOSE 8080

CMD ["/demo/gin-auth2-swagger-demo"]

# CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gin-auth2-swagger-demo .
# sudo docker build -t gin-auth2-swagger-demo .