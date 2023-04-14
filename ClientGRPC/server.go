package ClientGRPC

import (
	CGrpc "github.com/yan259128/alg_bcDB/Proto/client"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct{}

func (s *Server) Init() {
	go s.startService()
}
func (s *Server) startService() {
	// 初始化grpc对象
	grpcServer := grpc.NewServer() // 注册服务
	CGrpc.RegisterClientServiceServer(grpcServer, &Service{})
	// 创建监听
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Panic(err)
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Panic(err)
		}
	}(listen)

	// 绑定服务
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Panic(err)
	}
}
