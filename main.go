package main

import (
	"context"
	"fmt"
	// "github.com/vctt94/kamona/database"
	// api "github.com/vctt94/kamona/kauth"
	"log"
	// "net/http"
	// "strings"

	"os"
	// jwt "github.com/dgrijalva/jwt-go"
	// "github.com/grpc-ecosystem/grpc-gateway/runtime"
	// "github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/metadata"
)

// private type for Context keys
type contextKey int

const (
	clientIDKey contextKey = iota

	certPaths = "./cert/"
)

func credMatcher(headerName string) (mdName string, ok bool) {
	if headerName == "Token" || headerName == "Password" {
		return headerName, true
	}
	return "", false
}

// func authenticateClient(ctx context.Context, s *api.Server) (database.User, error) {
// 	if md, ok := metadata.FromIncomingContext(ctx); ok {
// 		clientToken := strings.Join(md["token"], "")
// 		fmt.Println("token", clientToken)
// 		token, _ := jwt.Parse(clientToken, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return "", fmt.Errorf("There was an error")
// 			}
// 			return []byte("secret"), nil
// 		})
// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if ok && token.Valid {
// 			var user database.User
// 			mapstructure.Decode(claims, &user)
// 			fmt.Println("got user", user)
// 			return s.Db.GetUserByEmail(user.Email)
// 		}
// 		return database.User{}, fmt.Errorf("Invalid authorization token")
// 	}
// 	return database.User{}, fmt.Errorf("missing authentication parameters")
// }

// unaryInterceptor calls authenticateClient with current context
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// s, ok := info.Server.(*api.Server)
	// if !ok {
	// 	return nil, fmt.Errorf("unable to cast server")
	// }
	// if info.FullMethod != "/kauth.Authentication/Signup" &&
	// 	info.FullMethod != "/kauth.Authentication/Login" {
	// 	user, err := authenticateClient(ctx, s)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	ctx = context.WithValue(ctx, "userid", user.ID)
	// 	return handler(ctx, req)
	// }
	return handler(ctx, req)
}

func startRESTServer(address, grpcAddress, certFile string) error {
	// ctx := context.Background()
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()
	// mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(credMatcher))
	// creds, err := credentials.NewClientTLSFromFile(certFile, "")
	// if err != nil {
	// 	return fmt.Errorf("could not load TLS certificate: %s", err)
	// }
	// // Setup the client gRPC options
	// opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	// // nopts := []grpc.DialOption{grpc.WithInsecure()}
	// // Register ping
	// err = api.RegisterPingHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	// if err != nil {
	// 	return fmt.Errorf("could not register service Ping: %s", err)
	// }

	// err = api.RegisterAuthenticationHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	// if err != nil {
	// 	return fmt.Errorf("could not register service Authentication: %s", err)
	// }

	// log.Printf("starting HTTP/1.1 REST server on %s", address)
	// http.ListenAndServe(address, mux)
	return nil
}

func run(ctx context.Context) error {
	// ToDO add this to a config file
	grpcAddress := fmt.Sprintf("%s:%d", "localhost", 7777)
	// restAddress := fmt.Sprintf("%s:%d", "localhost", 7778)
	certFile := certPaths + "server.crt"
	keyFile := certPaths + "server.key"

	err := startGRPCServer(grpcAddress, certFile, keyFile)
	if err != nil {
		log.Fatalf("failed to start gRPC server: %s", err)
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()
	// Run kamona until permanent failure or shutdown is requested.
	if err := run(ctx); err != nil && err != context.Canceled {
		os.Exit(1)
	}
}
