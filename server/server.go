package server

import (
	"github.com/yan259128/alg_bcDB/ClientGRPC"
	"github.com/yan259128/alg_bcDB/GRPC"
	"github.com/yan259128/alg_bcDB/blockchain/blockchain_data"
	"github.com/yan259128/alg_bcDB/blockchain/blockchain_table"
	"github.com/yan259128/alg_bcDB/cache"
	"github.com/yan259128/alg_bcDB/txpool"
	"github.com/yan259128/alg_bcDB/userManage"
	"sync"
)

// Server 定义一个区块链数据库服务
type Server struct {
	mutex      sync.Mutex
	manage     userManage.UserManage
	dataChain  blockchain_data.BlockChain
	tableChain blockchain_table.BlockChain
	TxPool     txpool.TxPool
	Cache      cache.Cache
	Grpc       GRPC.Server
	ClientGrpc ClientGRPC.Server
}

func (s *Server) Init() {
	s.manage.Init()
	s.dataChain.InitBlockChain()
	s.tableChain.InitBlockChain()
	s.Cache.Init(&s.dataChain, &s.tableChain)
	s.TxPool.Init(&s.dataChain, &s.tableChain, &s.Cache)
	s.Grpc.Init()
	s.ClientGrpc.Init()
}
