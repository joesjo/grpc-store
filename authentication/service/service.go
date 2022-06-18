package service

import (
	"context"
	"log"
	"net"

	"github.com/joesjo/grpc-store/authentication/database"
	pb "github.com/joesjo/grpc-store/authentication/protobuf"
	"google.golang.org/grpc"

	"github.com/go-playground/validator/v10"
	"github.com/joesjo/grpc-store/authentication/security"
	"golang.org/x/crypto/bcrypt"
)

const (
	port = ":8081"
)

type server struct {
	pb.UnimplementedAuthenticationServiceServer
}

type InvalidRequestError struct {
	message string
}

func (e *InvalidRequestError) Error() string {
	return "Invalid request: " + e.message
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Println("Creating user with username:", req.User.Username)
	validate := validator.New()
	var err error
	err = validate.Var(req.User.Username, "required,min=3,max=20")
	if err != nil {
		return nil, &InvalidRequestError{message: err.Error()}
	}
	err = validate.Var(req.User.Password, "required,min=8,max=20")
	if err != nil {
		return nil, &InvalidRequestError{message: err.Error()}
	}
	foundUser, err := database.FindUser(req.User.Username)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return nil, err
		}
	}
	if foundUser != nil {
		return nil, &InvalidRequestError{message: "user already exists"}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	err = database.CreateUser(req.User.Username, string(hashedPassword))
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{}, nil
}

func (s *server) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	if req.User.Username == "" {
		return nil, &InvalidRequestError{message: "user name is required"}
	}
	if req.User.Password == "" {
		return nil, &InvalidRequestError{message: "password is required"}
	}
	foundUser, err := database.FindUser(req.User.Username)
	if err != nil {
		return nil, err
	}
	if foundUser == nil {
		return nil, &InvalidRequestError{message: "user not found"}
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.User.Password))
	if err != nil {
		return nil, &InvalidRequestError{message: "invalid password"}
	}
	token, err := security.CreateToken(foundUser.Username)
	if err != nil {
		return nil, err
	}
	return &pb.AuthenticateResponse{Token: token}, nil
}

func (s *server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	username, err := security.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	return &pb.ValidateTokenResponse{Username: username}, nil
}

func Start() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthenticationServiceServer(s, &server{})
	log.Printf("Starting authentication server on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
