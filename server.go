package grpctmpl

import (
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type server struct {
	mux   *runtime.ServeMux
	port  int
	grpcS *grpc.Server
}

func New(port int) *server {
	//var opts = []grpc.ServerOption{
	//	grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	//	grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	//}
	//grpcServer := grpc.NewServer(opts...)

	newServer := server{
		port: port,
		mux:  runtime.NewServeMux(),
	}
	return &newServer
}
func (s *server) Register(grpcRegister func(srv grpc.ServiceRegistrar), httpRegister func(mux *runtime.ServeMux, endpoint string) (err error)) error {
	grpcS := grpc.NewServer([]grpc.ServerOption{}...)
	grpcRegister(grpcS)
	s.grpcS = grpcS
	//grpc_prometheus.Register(s.grpcS)

	return httpRegister(s.mux, fmt.Sprintf(":%d", s.port))
}
func (s *server) Serve() error {
	//if err := s.mux.HandlePath("GET", "/test/{name}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	//	w.Write([]byte("hello " + pathParams["name"]))
	//}); err != nil {
	//	return err
	//}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		panic(err)
	}

	m := cmux.New(l)

	grpcL := m.Match(cmux.HTTP2())
	httpL := m.Match(cmux.Any())

	httpS := &http.Server{
		Handler: s.mux,
	}
	go s.grpcS.Serve(grpcL)

	go httpS.Serve(httpL)


	log.Printf("server listen at port: %d\n", s.port)
	return m.Serve()
}
