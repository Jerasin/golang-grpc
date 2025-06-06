package middleware

import (
	"auth-svc/pkg/utils"
	"context"
	"log"
	"slices"

	"github.com/goforj/godump"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthorizationInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	jwt utils.JWTWrapper,
	options *[]CustomParameterMiddleware,
) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	if slices.ContainsFunc(*options, func(u CustomParameterMiddleware) bool {
		return u.SkipPath != info.FullMethod
	}) {
		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
		}

		godump.Dump("gRPC called values: %v", values)
		_, validateTokenErr := jwt.ValidateToken(values[0], "user")
		if validateTokenErr != nil {
			return nil, status.Errorf(codes.Unauthenticated, "validate token failed")
		}
	}

	res, err := handler(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return res, nil
}
