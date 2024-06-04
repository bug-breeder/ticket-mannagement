package service

import (
	"context"
	"time"

	"go.tekoapis.com/tekone/app/warehouse/iam_service/api"
	health "go.tekoapis.com/tekone/library/grpc/health"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// "go.tekoapis.com/tekone/app/warehouse/iam_service/config"
	"go.tekoapis.com/tekone/app/warehouse/iam_service/custom_constants"
	"go.tekoapis.com/tekone/app/warehouse/iam_service/internal/store"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type Service struct {
	log logr.Logger
	// more connector here
	store store.StoreQuerier
	// embedded unimplemented service server
	health.UnimplementedHealthCheckServiceServer
	api.UnimplementedIamServiceServer
}

func NewService(
	logger logr.Logger,
	store store.StoreQuerier,
) *Service {

	return &Service{
		log:   logger,
		store: store,
		// more here

	}
}

func (s *Service) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, err
	}

	// // username must be unique
	// _, err = s.store.GetUserByUsername(ctx, req.Username)
	// if err == nil {
	// 	return nil, status.Errorf(codes.AlreadyExists, "username already exists")
	// }

	// // role must be either manager' or 'employee'
	// if req.Role != "manager" && req.Role != "employee" {
	// 	return nil, status.Errorf(codes.InvalidArgument, "role must be either 'manager' or 'employee'")
	// }

	// // gender must be male, female or other
	// if req.Gender != "male" && req.Gender != "female" && req.Gender != "other" {
	// 	return nil, status.Errorf(codes.InvalidArgument, "gender must be either 'male', 'female', or 'other'")
	// }

	userID, err := s.store.CreateUser(ctx, store.CreateUserParams{
		FullName:  req.FullName,
		Username:  req.Username,
		Gender:    req.Gender,
		BirthDate: birthDate,
		Password:  req.Password,
		Role:      req.Role,
	})
	if err != nil {
		return nil, err
	}

	return &api.CreateUserResponse{
		Message: "User created successfully",
		UserId:  userID.String(),
	}, nil
}

func (s *Service) GetToken(ctx context.Context, req *api.GetTokenRequest) (*api.GetTokenResponse, error) {
	user, err := s.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// Verify the password

	if user.Password != req.Password {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID.String(),
		"role":    user.Role,
	})
	tokenString, err := token.SignedString([]byte(custom_constants.JWTSecretKey))
	if err != nil {
		return nil, err
	}

	return &api.GetTokenResponse{
		Token: tokenString,
	}, nil
}

func (s *Service) GetUserInfo(ctx context.Context, req *api.GetUserInfoRequest) (*api.GetUserInfoResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &api.GetUserInfoResponse{
		UserId:    user.UserID.String(),
		FullName:  user.FullName,
		Username:  user.Username,
		Gender:    user.Gender,
		BirthDate: user.BirthDate.String(),
		Role:      user.Role,
	}, nil
}

func (s *Service) Close(ctx context.Context) {
	// log "closing service"
	s.log.Info("closing service")
	s.store.Close()
}

func (s *Service) Ping() error {
	err := s.store.Ping()
	return err
}
