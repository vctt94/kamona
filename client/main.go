package main

import (
	"context"
	api "kamona/kauth"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Authentication struct {
	Login    string
	Password string
}

// GetRequestMetadata gets the current request metadata
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security
func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

func main() {
	var conn *grpc.ClientConn
	serverCert := "/Users/fernandoabolafio/go/src/grpc_play/cert/server.crt"
	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile(
		serverCert,
		"localhost",
	)
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	auth := Authentication{
		Login:    "john",
		Password: "doe",
	}

	conn, err = grpc.Dial(":7777", grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	defer conn.Close()

	c := api.NewPingClient(conn)

	response, err := c.SayHello(context.Background(), &api.PingMessage{Greeting: "foo"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Greeting)
}
