package kauth

import (
	"context"
	"fmt"
	"github.com/getset0/kamona/database"
	"log"
	"sync"

	jwt "github.com/dgrijalva/jwt-go"
)

// Server represents the gRPC server
type Server struct {
	sync.RWMutex

	Db *database.Database
}

// NewServer returns a new server
func NewServer(db *database.Database) Server {
	return Server{
		Db: db,
	}
}

// SayHello generates response to a Ping request
func (s *Server) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	log.Printf("Receive message %s", in.Greeting)
	return &PingMessage{Greeting: "bar"}, nil
}

// Signup registers a new user
func (s *Server) Signup(ctx context.Context, in *SingupInput) (*SignupOutput, error) {

	s.Db.NewUser(database.User{
		Email:    in.Email,
		Password: in.Password,
	})
	log.Printf("User registered %v", in.Email)

	return &SignupOutput{
		Success: true,
	}, nil
}

// Login authenticate an user
func (s *Server) Login(ctx context.Context, in *LoginInput) (*LoginOutput, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    in.Email,
		"password": in.Password,
	})

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	return &LoginOutput{
		Token:   tokenString,
		Success: true,
	}, nil
}

func (s *Server) Me(ctx context.Context, in *MeInput) (*User, error) {
	userID := ctx.Value("userid").(string)
	user, err := s.Db.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:    user.ID,
		Email: user.Email,
	}, nil
}
