package services

import (
	"auth-svc/pkg/models"
	"auth-svc/pkg/pb"
	"auth-svc/pkg/repositories"
	"auth-svc/pkg/utils"
	"context"

	"github.com/goforj/godump"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	UserRepo repositories.BaseRepository[*models.User]
	Jwt      *utils.JWTWrapper
}

func (s *Server) Test(_ context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	return &pb.TestResponse{
		Message: "Hello " + req.Name,
	}, nil
}

func (s *Server) RegisterUser(c context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := utils.Validate(&user); err != nil {
		return nil, err
	}
	_, err := s.UserRepo.Transaction(c, func(ctx context.Context) (any, error) {
		var err error
		if _, err = s.UserRepo.IsExist(ctx, map[string]any{"email": req.Email}); err != nil {
			return nil, err
		}

		_, err = s.UserRepo.Save(ctx, &models.User{
			Email:    req.Email,
			Password: utils.HashPassword(req.Password),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to save user: %v", err)
		}
		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Status: 200,
	}, nil
}

func (s *Server) LoginUser(c context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := utils.Validate(&user); err != nil {
		return nil, err
	}
	jwtToken, err := s.UserRepo.Transaction(c, func(ctx context.Context) (any, error) {
		var err error
		var res *models.User
		if res, err = s.UserRepo.FindOne(ctx, map[string]any{"email": req.Email}); err != nil {
			return nil, status.Errorf(codes.NotFound, "User not found: %v", err)
		}

		if !utils.CheckPasswordHash(req.Password, res.Password) {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid password")
		}

		token, _ := s.Jwt.GenerateToken(*res, "user")

		godump.Dump("token", token)
		return token, nil
	})

	if err != nil {
		return nil, err
	}

	val, ok := jwtToken.(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to generate JWT token")
	}

	godump.Dump(val)

	return &pb.LoginResponse{
		Status:      200,
		AccessToken: val,
	}, nil
}
