package server

import (
	"fmt"
	"github.com/yan259128/alg_bcDB/GRPC"
	"github.com/yan259128/alg_bcDB/blockchain/blockchain_data"
	"github.com/yan259128/alg_bcDB/blockchain/blockchain_table"
	"time"
)

func (s *Server) RegisterAccount(username, password string) bool {

	err := s.manage.Register(username, password)
	if err != nil {
		fmt.Printf("注册用户失败; %s\n", err)
		return false
	}
	fmt.Printf("注册用户成功; \n")
	return true
}

func (s *Server) SignInAccount(username, password string) bool {

	err := s.manage.SignIn(username, password)
	if err != nil {
		fmt.Printf("登录失败; %s\n", err)
		return false
	}
	fmt.Printf("登录成功; \n")
	return true
}

func (s *Server) CheckAccount(UID string) string {

	a, err := s.manage.ViewAccount(UID)
	if err != nil {
		fmt.Println("查看用户信息失败 ;", err)
		return ""
	}
	fmt.Printf("address %s\n", a.Address)
	return a.Address
}

func (s *Server) SignOutAccount(ID string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	err := s.manage.SignOut(ID)
	if err != nil {
		fmt.Printf("退出失败; %s\n", err)
		return false
	}
	fmt.Printf("用户退出成功\n")
	return true
}

func (s *Server) Table(UID string, table string, permissionTable []string) bool {

	// 登录状态检查
	a, err := s.manage.ViewAccount(UID)
	if err != nil {
		fmt.Println("用户未登录")
		return false
	}
	ad := a.Address

	// 权限检查
	permission, err := s.Cache.CheckPermission(ad, table)
	if err != nil {
		fmt.Printf("创建一个新的表, %s\n", table)
		goto alter
	} else {
		if permission < "4" {
			fmt.Println("没有对表的修改权限")
			return false
		}
	}

alter:
	var tx blockchain_table.Transaction
	tx.Init(table, permissionTable, a.UserName, a.PublicKey, a.PrivateKey)

	if blockchain_table.VerifyTransaction(tx) {

		err := s.TxPool.TxTableIN(tx)
		if err != nil {
			fmt.Println(err)
		}
		// 交易的广播
		GRPC.SubmitTableTransaction(tx)
		fmt.Printf("权限表信息校验成功.\n")
		return true
	} else {
		fmt.Printf("权限表信息校验失败.\n")
		return false
	}
}

func (s *Server) Put(UID string, key, value, table string) bool {

	// 登录状态检查
	a, err := s.manage.ViewAccount(UID)
	if err != nil {
		fmt.Println("用户未登录")
		return false
	}
	ad := a.Address

	// 权限检查
	permission, err := s.Cache.CheckPermission(ad, table)
	if err != nil {
		fmt.Printf("不存在这个表; %s\n", table)
		return false
	} else {
		if permission < "3" {
			fmt.Println("没有对表的写入权限")
			return false
		}
	}

	// 1. 创建交易（这是本地的处理，如果是其他节点的交易直接校验并添加数据到交易池）
	var tx blockchain_data.Transaction
	tx.Init(table, key, value, a.UserName, a.PublicKey, a.PrivateKey)

	// 2.验证交易
	// 检验并添加到交易池
	if blockchain_data.VerifyTransaction(tx) {
		fmt.Printf("数据写入交易校验成功.\n")

		err := s.TxPool.TxDataIN(tx)
		if err != nil {
			fmt.Println(err)
		}
		// 交易的广播
		GRPC.SubmitDataTransaction(tx)
		return true
	} else {
		fmt.Printf("数据写入交易校验失败.\n")
		return false
	}
}

func (s *Server) Get(UID, key, table string) string {
	start := time.Now() // 获取当前时间

	// 登录状态检查
	a, err := s.manage.ViewAccount(UID)
	if err != nil {
		fmt.Println("用户未登录")
		return ""
	}
	ad := a.Address

	// 权限检查
	permission, err := s.Cache.CheckPermission(ad, table)
	if err != nil {
		fmt.Printf("不存在这个表; %s\n", table)
		return ""
	} else {
		if permission < "1" {
			fmt.Println("没有对表的查看权限")
			return ""
		}
	}
	// 验证权限

	tx, err := s.Cache.GetOneValue(table+"-QAQ-"+key, table)

	if err != nil {
		fmt.Printf("没有找到该数据\n")
		return ""
	}
	fmt.Printf("key : %s    ", tx.Key)
	fmt.Printf("value: %s    ", tx.Value)
	fmt.Printf("possessor: %s    ", tx.Possessor)
	fmt.Printf("alterTime : %v\n", time.Unix(tx.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
	fmt.Printf("(查询结束)\n")

	elapsed := time.Since(start)
	fmt.Println("该查找执行完成耗时：", elapsed)
	return tx.Value
}

func (s *Server) GetHistory(UID, key, table string) {
	start := time.Now() // 获取当前时间

	// 登录状态检查
	a, err := s.manage.ViewAccount(UID)
	if err != nil {
		fmt.Println("用户未登录")
		return
	}
	ad := a.Address

	// 权限检查
	permission, err := s.Cache.CheckPermission(ad, table)
	if err != nil {
		fmt.Printf("不存在这个表; %s\n", table)
		return
	} else {
		if permission < "1" {
			fmt.Println("没有对表的查看权限")
			return
		}
	}

	txs, err := s.Cache.GetHistoryValue(table+"-QAQ-"+key, table)
	for _, tx := range txs {
		fmt.Printf("key : %s    ", tx.Key)
		fmt.Printf("value: %s    ", tx.Value)
		fmt.Printf("possessor: %s    ", tx.Possessor)
		fmt.Printf("alterTime : %v\n", time.Unix(tx.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
	}
	elapsed := time.Since(start)
	fmt.Println("该查找执行完成耗时：", elapsed)
}

func (s *Server) MyTable(UID string) []string {

	var tabelNames []string
	// 登录状态检查
	a, err := s.manage.ViewAccount(UID)
	if err != nil {
		fmt.Println("用户未登录")
		return []string{}
	}
	ad := a.Address

	tables, err := s.Cache.MyTables(ad)
	if err != nil {
		fmt.Println("没有任何表的权限")
		return []string{}
	}

	for tn, pl := range tables {
		fmt.Printf("table : %s, privilege level : %s\n", tn, pl)
		tabelNames = append(tabelNames, tn)
	}
	fmt.Println()
	return tabelNames
}
