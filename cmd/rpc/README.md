# shell

```(shell)
DIR=./proto/echo
protoc -I $DIR $DIR/*.proto --go_out=plugins=grpc:$DIR
```

```(shell)
openssl genrsa -out server.key 2048
openssl req -x509 -key server.key -out server.pem
```
