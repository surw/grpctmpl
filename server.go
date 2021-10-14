package grpctmpl

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"go.elastic.co/apm/module/apmgrpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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
		mux: runtime.NewServeMux(
			runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, message proto.Message) error {
				//headers := w.Header()
				//if location, ok := headers[context_helper.HeaderLocationKey]; ok {
				//	w.Header().Set(context_helper.HeaderLocationKey, location[0])
				//	w.WriteHeader(http.StatusFound)
				//}
				//if len(headers.Get(context_helper.HeaderContentDispositionKey)) > 0 {
				//	w.Header().Set("content-transfer-encoding", "binary")
				//	w.Header().Set("content-type", "application/force-download")
				//	w.Header().Add("content-type", "application/octet-stream")
				//}

				return nil
			}),
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					UseEnumNumbers:  true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			}),
			runtime.WithMarshalerOption("application/json+pretty", &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					Indent:          "  ",
					UseProtoNames:   true,
					UseEnumNumbers:  true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			}),
			//runtime.WithMarshalerOption(grpc_marshaller.TextCsv, grpc_marshaller.ChunkMarshaller{}),
			//runtime.WithMarshalerOption(grpc_marshaller.ImageJpeg, grpc_marshaller.ChunkMarshaller{}),
			//runtime.WithIncomingHeaderMatcher(IncomingHeaderMatcher),
			//runtime.WithProtoErrorHandler(DefaultProtoErrorHandler),
			//runtime.WithOutgoingHeaderMatcher(OutgoingHeaderMatcherForSuccess),
			//runtime.WithMetadata(AppendRequestMetadata),
		),
	}
	return &newServer
}
func (s *server) Register(grpcRegister func(srv grpc.ServiceRegistrar), httpRegister func(mux *runtime.ServeMux, endpoint string) (err error)) error {
	s.grpcS = grpc.NewServer(grpc.UnaryInterceptor(
		apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
	))
	grpcRegister(s.grpcS)
	reflection.Register(s.grpcS)
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

	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))

	httpL := m.Match(cmux.Any())

	//cer, err := tls.LoadX509KeyPair("cert", "private")
	//if err != nil {
	//	panic(err)
	//}
	//
	//config := &tls.Config{Certificates: []tls.Certificate{cer}}

	mux := http.NewServeMux()
	mux.Handle("/", s.mux)

	fs := http.FileServer(http.Dir("/swagger"))
	mux.Handle("/help/", http.StripPrefix("/help", fs))

	httpS := &http.Server{
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	//go httpS.ServeTLS(httpL, "cert", "private")
	go httpS.Serve(httpL)
	go s.grpcS.Serve(grpcL)

	log.Printf("server listen at port: %d\n", s.port)
	return m.Serve()
}
