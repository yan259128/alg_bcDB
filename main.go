package main

import (
	"fmt"
	"github.com/yan259128/alg_bcDB/ClientGRPC"
	"github.com/yan259128/alg_bcDB/Cluster"
	"github.com/yan259128/alg_bcDB/Raft"
	"github.com/yan259128/alg_bcDB/algorand"
	"github.com/yan259128/alg_bcDB/blockqueue"
	"github.com/yan259128/alg_bcDB/client"
	"github.com/yan259128/alg_bcDB/server"
	"github.com/yan259128/alg_bcDB/serverExec"
	"log"
	"os"
)

func main() {

	blockqueue.LocalDataBlockQueue = blockqueue.Init()
	blockqueue.LocalTableBlockQueue = blockqueue.Init()
	//blockqueue.AlgToBC = blockqueue.Init()
	//blockqueue.BCToAlg = blockqueue.Init()
	s := new(server.Server)
	serverExec.RPCs = s
	client.Cserver = s
	ClientGRPC.RPCs = s
	s.Init()

	// 判断集群文件是否存在，如果不存在则直接执行cmd程序
	_, err := os.Stat("./ClusterInfo")
	if os.IsNotExist(err) {
		go s.Command()
		go client.StartClient()
		//serverExec.ServerStart()

	} else {
		// 启动Raft, 读取集群文件
		cluster, err := Cluster.LoadClusterFile("./ClusterInfo")
		if err != nil {
			log.Panic(err)
		}

		Cluster.LocalNode = cluster
		for _, node := range cluster.Node {
			fmt.Println(node.IP, node.Port)
		}

		go Cluster.Server()
		// 集群中节点的更新
		go cluster.UpdateClusterFile()

		go Raft.Start(&s.TxPool)

		alg := algorand.NewAlgorand()
		//go Raft.Start(&s.TxPool)
		alg.Start()

		go s.Command()
		go client.StartClient()
		//serverExec.ServerStart()
	}

}
