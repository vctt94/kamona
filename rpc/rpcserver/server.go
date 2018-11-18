package rpcserver

import (
	"context"
	"fmt"
	pb "github.com/getset0/kamona/rpc/kamonarpc"
	"github.com/vctt94/kamona/database"
	"google.golang.org/grpc"
	"log"
	"sync"

	jwt "github.com/dgrijalva/jwt-go"
)

// Public API version constants
const (
	semverString = "0.0.0"
	semverMajor  = 0
	semverMinor  = 0
	semverPatch  = 0
)

// versionServer provides RPC clients with the ability to query the RPC server
// version.
type versionServer struct{}

type authServer struct {
	sync.RWMutex

	Db *database.Database
}

var (
	versionService versionServer
	authService    authServer
)

// RegisterServices registers implementations of each gRPC service and registers
// it with the server.
func RegisterServices(server *grpc.Server) {
	pb.RegisterVersionServiceServer(server, &versionService)
	pb.RegisterAuthServiceServer(server, &authService)
}

var serviceMap = map[string]interface{}{
	"kamonarpc.VersionService": &versionService,
	"kamonarpc.AuthService":    &authService,
}

// Version Service
func (*versionServer) Version(ctx context.Context, req *pb.VersionRequest) (*pb.VersionResponse, error) {
	return &pb.VersionResponse{
		VersionString: semverString,
		Major:         semverMajor,
		Minor:         semverMinor,
		Patch:         semverPatch,
	}, nil
}

// Auth Service

// Signup registers a new user
func (s *authServer) Signup(ctx context.Context, in *pb.SingupInput) (*pb.SignupOutput, error) {

	s.Db.NewUser(database.User{
		Email:    in.Email,
		Password: in.Password,
	})
	log.Printf("User registered %v", in.Email)

	return &pb.SignupOutput{
		Success: true,
	}, nil
}

// Login authenticate an user
func (s *authServer) Login(ctx context.Context, in *pb.LoginInput) (*pb.LoginOutput, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    in.Email,
		"password": in.Password,
	})

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	return &pb.LoginOutput{
		Token:   tokenString,
		Success: true,
	}, nil
}

func (s *authServer) Me(ctx context.Context, in *pb.MeInput) (*pb.User, error) {
	userID := ctx.Value("userid").(string)
	user, err := s.Db.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:    user.ID,
		Email: user.Email,
	}, nil
}
