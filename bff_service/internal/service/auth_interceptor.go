package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.tekoapis.com/tekone/app/warehouse/bff_service/custom_constants"
)

func parseToken(tokenString string) (jwt.MapClaims, error) {
	jwtKey := custom_constants.JWTSecretKey

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func (s *Service) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Apply middleware only to specific method
	if info.FullMethod != "/tekone.app.warehouse.bff_service.api.BffService/UpdateTicketStatus" {
		log.Printf("skip auth middleware for method %s", info.FullMethod)
		return handler(ctx, req)
	}

	// log "starting auth middleware"
	log.Printf("starting auth middleware for method %s", info.FullMethod)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing context metadata")
	}

	authorization := md["authorization"]
	if len(authorization) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
	}

	tokenString := authorization[0]
	// Remove "Bearer " prefix if present
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	claims, err := parseToken(tokenString)
	if err != nil {
		log.Printf("invalid token: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	role, ok := claims["role"].(string)
	if !ok || role != "manager" {
		return nil, status.Errorf(codes.PermissionDenied, "only managers can approve/reject tickets")
	}
	log.Printf("auth middleware passed for method %s", info.FullMethod)

	return handler(ctx, req)
}
