package main

import (
	"context"
	"fmt"
	"kamona/database"
	api "kamona/kauth"
	"log"
	"net"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// private type for Context keys
type contextKey int

const (
	clientIDKey contextKey = iota

	certPaths = "/Users/fernandoabolafio/go/src/kamona/cert/"
)

func credMatcher(headerName string) (mdName string, ok bool) {
	if headerName == "Token" || headerName == "Password" {
		return headerName, true
	}
	return "", false
}

func authenticateClient(ctx context.Context, s *api.Server) (database.User, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientToken := strings.Join(md["token"], "")
		fmt.Println("token", clientToken)
		token, _ := jwt.Parse(clientToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return "", fmt.Errorf("There was an error")
			}
			return []byte("secret"), nil
		})
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			var user database.User
			mapstructure.Decode(claims, &user)
			fmt.Println("got user", user)
			return s.Db.GetUserByEmail(user.Email)
		}
		return database.User{}, fmt.Errorf("Invalid authorization token")
	}
	return database.User{}, fmt.Errorf("missing authentication parameters")
}

// unaryInterceptor calls authenticateClient with current context
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	s, ok := info.Server.(*api.Server)
	if !ok {
		return nil, fmt.Errorf("unable to cast server")
	}
	if info.FullMethod != "/kauth.Authentication/Signup" &&
		info.FullMethod != "/kauth.Authentication/Login" {
		user, err := authenticateClient(ctx, s)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, "userid", user.ID)
		return handler(ctx, req)
	}
	return handler(ctx, req)
}

func startGRPCServer(address, certFile, keyFile string) error {
	// create a listener on TCP port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	// create a server instance
	db := database.NewDatabase()
	s := api.NewServer(db)

	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("could not load TLS keys: %s", err)
	}
	// Create an array of gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor)}
	// create a gRPC server object
	grpcServer := grpc.NewServer(opts...)

	// attach the Ping service to the server
	api.RegisterPingServer(grpcServer, &s)
	api.RegisterAuthenticationServer(grpcServer, &s)

	// start the server
	log.Printf("starting HTTP/2 gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}
	return nil
}

func startRESTServer(address, grpcAddress, certFile string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(credMatcher))
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return fmt.Errorf("could not load TLS certificate: %s", err)
	}
	// Setup the client gRPC options
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	// nopts := []grpc.DialOption{grpc.WithInsecure()}
	// Register ping
	err = api.RegisterPingHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return fmt.Errorf("could not register service Ping: %s", err)
	}

	err = api.RegisterAuthenticationHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return fmt.Errorf("could not register service Authentication: %s", err)
	}

	log.Printf("starting HTTP/1.1 REST server on %s", address)
	http.ListenAndServe(address, mux)
	return nil
}

func main() {
	grpcAddress := fmt.Sprintf("%s:%d", "localhost", 7777)
	restAddress := fmt.Sprintf("%s:%d", "localhost", 7778)
	certFile := certPaths + "server.crt"
	keyFile := certPaths + "server.key"

	go func() {
		err := startGRPCServer(grpcAddress, certFile, keyFile)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	// fire the REST server in a goroutine
	go func() {
		err := startRESTServer(restAddress, grpcAddress, certFile)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	// infinite loop
	log.Printf("Entering infinite loop")
	select {}
}
