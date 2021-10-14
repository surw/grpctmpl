package server_interceptors

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmlogrus"

	"google.golang.org/grpc"
)

func FallbackInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}
func LogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		logger := log.WithFields(apmlogrus.TraceContext(ctx))
		logger.Info("req", req, "method", info.FullMethod)
		resp, err = handler(ctx, req)
		return resp, err
	}
}

var log = &logrus.Logger{
	Out:   os.Stderr,
	Hooks: make(logrus.LevelHooks),
	Level: logrus.DebugLevel,
	Formatter: &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "log.level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "function.name", // non-ECS
		},
	},
}
