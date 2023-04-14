package ClientGRPC

import (
	"context"
	CGrpc "github.com/yan259128/alg_bcDB/Proto/client"
	"github.com/yan259128/alg_bcDB/server"
	"strings"
)

// Service GRPC接口的实现(客户端)
type Service struct{}

var RPCs *server.Server

func (s Service) Login(ctx context.Context, mes *CGrpc.LoginMes) (info *CGrpc.VerifyInfo, err error) {
	info = &CGrpc.VerifyInfo{}
	info.Info = "登录成功"
	info.Status = true
	if mes.UserName == "" || mes.Password == "" {
		info.Info = "登录失败"
		info.Status = false
	} else {
		if !RPCs.SignInAccount(mes.UserName, mes.Password) {
			info.Info = "登录失败"
			info.Status = false
		}
	}
	return info, nil
}

func (s Service) CreatTable(ctx context.Context, table *CGrpc.TableInfo) (info *CGrpc.VerifyInfo, err error) {
	info = &CGrpc.VerifyInfo{}
	info.Info = "创建成功"
	info.Status = true
	if table.TableName == "" || table.Permission == "" {
		info.Info = "创建失败"
		info.Status = false
	} else {
		permissionTable := strings.Split(table.Permission, ",")
		isCreatetabelSuccess := RPCs.Table(table.Uid, table.TableName, permissionTable)
		if !isCreatetabelSuccess {
			info.Info = "创建失败"
			info.Status = false
		}
	}
	return info, nil
}

func (s Service) Read(ctx context.Context, read *CGrpc.ReadInfo) (info *CGrpc.VerifyInfo, err error) {
	info = &CGrpc.VerifyInfo{}
	info.Info = "读取成功"
	info.Status = true
	if read.Uid == "" || read.Key == "" || read.TableName == "" {
		info.Info = "读取失败"
		info.Status = false
	} else {
		getValue := RPCs.Get(read.Uid, read.Key, read.TableName)
		if getValue == "" {
			info.Info = "读取失败"
			info.Status = false
		}
	}
	return info, nil
}

func (s Service) Write(ctx context.Context, write *CGrpc.WriteInfo) (info *CGrpc.VerifyInfo, err error) {
	info = &CGrpc.VerifyInfo{}
	info.Info = "写入成功"
	info.Status = true
	if write.Uid == "" || write.Key == "" || write.Value == "" || write.TableName == "" {
		info.Info = "写入失败"
		info.Status = false
	} else {
		isPutSuccess := RPCs.Put(write.Uid, write.Key, write.Value, write.TableName)
		if !isPutSuccess {
			info.Info = "写入失败"
			info.Status = false
		}
	}
	return info, nil
}
