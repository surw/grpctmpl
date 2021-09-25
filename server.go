package grpctmpl

import (
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
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
	return registerFunc(s.mux, fmt.Sprintf(":%d", s.port))
}
