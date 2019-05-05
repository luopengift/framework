package main

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/luopengift/framework/cmd/rpc/proto/echo"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

const (
	address = "localhost:50051"
)

// fetchToken simulates a token lookup and omits the details of proper token
// acquisition. For examples of how to acquire an OAuth2 token, see:
// https://godoc.org/golang.org/x/oauth2
func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}

func main() {

	perRPC := oauth.NewOauthAccess(fetchToken())
	opts := []grpc.DialOption{
		// grpc.WithInsecure(),
		// In addition to the following grpc.DialOption, callers may also use
		// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
		// itself.
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		grpc.WithPerRPCCredentials(perRPC),
		// oauth.NewOauthAccess requires the configuration of transport
		// credentials.
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
		),
	}

	//建立链接
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := echo.NewEchoClient(conn)

	msg := time.Now().String()

	// 1秒的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Echo(ctx, &echo.Request{Msg: msg})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Echo: %s", r.Msg)
}
