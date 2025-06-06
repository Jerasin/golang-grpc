package middleware

import (
	"auth-svc/pkg/models"
	"auth-svc/pkg/repositories"
	"context"
	"encoding/json"
	"net"
	"slices"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	auditLogRepo repositories.BaseRepository[*models.AuditLog],
	options *[]CustomParameterMiddleware,
) (any, error) {
	var jsonRes string
	var handlerErr error
	start := time.Now()
	res, handlerErr := handler(ctx, req)
	duration := time.Since(start)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	st, ok := status.FromError(handlerErr)
	if !ok {
		st = status.New(codes.Unknown, handlerErr.Error())
	}

	if res != "" {
		jsonResData, err := json.Marshal(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error")
		}
		jsonRes = string(jsonResData)
	}

	if slices.ContainsFunc(*options, func(u CustomParameterMiddleware) bool {
		return u.SkipPath != info.FullMethod
	}) {
		var traceID []string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			traceID = md["trace-id"]
		}

		p, ok := peer.FromContext(ctx)
		if !ok || p.Addr == net.Addr(nil) {
			return nil, status.Errorf(codes.Aborted, "Aborted")
		}

		addrArr := strings.Split(p.Addr.String(), ":")
		if len(addrArr) < 2 {
			status.Error(codes.FailedPrecondition, "addrArr error")
		}

		auditLog := &models.AuditLog{
			Method:   info.FullMethod,
			Request:  string(jsonData),
			Response: jsonRes,
			ClientIP: addrArr[0],
			Status:   st.Code().String(),
			Error:    st.Message(),
			Duration: duration.Milliseconds(),
			Protocol: "grpc",
			TraceID:  traceID,
		}
		if len(traceID) == 0 {
			auditLog.Status = codes.InvalidArgument.String()
			auditLog.Error = "traceID is required"
			handlerErr = status.Error(codes.InvalidArgument, "traceID is required")

		}

		_, err := auditLogRepo.Save(ctx, auditLog)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if handlerErr != nil {
			return nil, handlerErr
		}

	}

	return res, nil
}
