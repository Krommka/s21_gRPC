package server

import (
	"fmt"
	"net"

	"Go_Team00.ID_376234-Team_TL_barievel/api/gen/pb"
	"Go_Team00.ID_376234-Team_TL_barievel/configs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	lis        net.Listener
}

func NewServer(cfg *configs.Config) *Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC.Port))
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(cfg.GRPC.MaxMessageSizeMB*1024*1024),
		grpc.MaxSendMsgSize(cfg.GRPC.MaxMessageSizeMB*1024*1024),
		grpc.ConnectionTimeout(cfg.GRPC.ConnectionTimeout),
	)
	pb.RegisterFrequencyServiceServer(server, newHandler())
	reflection.Register(server)

	return &Server{
		grpcServer: server,
		lis:        lis,
	}
}

func (s *Server) Serve() error {
	return s.grpcServer.Serve(s.lis)
}

func (s *Server) Shutdown() {
	s.grpcServer.GracefulStop()
}
