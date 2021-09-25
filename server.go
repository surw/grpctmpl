package grpctmpl

import (
	"fmt"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type server struct {
	mux  *runtime.ServeMux
	port int
}

func New(port int) *server {
	//var opts = []grpc.ServerOption{
	//	grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	//	grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	//}
	//grpcServer := grpc.NewServer(opts...)

	mux := runtime.NewServeMux()
	newServer := server{
		port: port,
		mux: mux,
	}
	return &newServer
}
func (s *server) RegisterHTTPGW(registerFunc func(mux *runtime.ServeMux, endpoint string) (err error)) error {
	grpcServer := grpc.NewServer()
	grpc_prometheus.Register(grpcServer)

	return registerFunc(s.mux, fmt.Sprintf(":%d", s.port))
}
