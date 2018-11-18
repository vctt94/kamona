package main

import (
	"context"
	"log"

	api "github.com/getset0/kamona/rpc/kamonarpc"

	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// this is just a test client
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

func client() {
	certificateFile := "./cert/server.crt"

	creds, err := credentials.NewClientTLSFromFile(certificateFile, "localhost")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := grpc.Dial("localhost:7777", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	c := api.NewVersionServiceClient(conn)
	versionRequest := &api.VersionRequest{}

	response, err := c.Version(context.Background(), versionRequest)

	if err != nil {
		log.Fatalf("Error when calling request: %s", err)
	}
	log.Printf("Response from server: %+v", response)
}
