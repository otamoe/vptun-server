package server

import (
	"errors"

	pb "github.com/otamoe/vptun-pb"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type (
	GrpcClientShell struct {
		Id         string
		grpcClient *GrpcClient
		close      *atomic.Bool
	}
)

func (grpcClientShell *GrpcClientShell) Close(status int32, werr error) (rerr error) {
	grpcClientShell.grpcClient.mux.Lock()
	delete(grpcClientShell.grpcClient.shells, grpcClientShell.Id)
	grpcClientShell.grpcClient.mux.Unlock()

	// 已执行过了
	if !grpcClientShell.close.CAS(false, true) {
		return
	}

	// 写入到数据库
	_, rerr = grpcClientShell.grpcClient.grpcHandler.clientShellSystem.Save(grpcClientShell.grpcClient.Id, grpcClientShell.Id, func(clientShell *ClientShell) (rClientShell *ClientShell, err error) {
		if status < 0 {
			status = 1
		}
		output := clientShell.Output
		if werr != nil {
			if len(output) != 0 {
				output += "\n"
			}
			output += werr.Error()
		}
		rClientShell = clientShell.WithStatus(status).WithOutput(output)
		return
	})

	fields := []zap.Field{
		zap.String("clientId", grpcClientShell.grpcClient.Id),
		zap.String("shellId", grpcClientShell.Id),
		zap.Int32("status", status),
	}

	if rerr != nil {
		grpcClientShell.grpcClient.logger("shell-close", false, rerr, fields...)
	} else {
		grpcClientShell.grpcClient.logger("shell-close", false, werr, fields...)
	}
	return
}

func (grpcClient *GrpcClient) OnShell(shellRequest *pb.ShellRequest) (err error) {
	if shellRequest.Id == "" {
		return
	}

	grpcClientShell := grpcClient.Shell(shellRequest.Id)

	// 没有返回取消
	if grpcClientShell == nil {
		err = grpcClient.Response(&pb.StreamResponse{
			Shell: &pb.ShellResponse{Id: shellRequest.Id},
		})
		return
	}

	// 关闭
	if shellRequest.Status != -1 {
		var werr error
		if len(shellRequest.Data) != 0 {
			werr = errors.New(string(shellRequest.Data))
		}
		err = grpcClientShell.Close(shellRequest.Status, werr)
		return
	}

	// 错误
	defer func() {
		if err != nil {
			go grpcClientShell.Close(1, err)
		}
	}()

	var clientShell *ClientShell
	if clientShell, err = grpcClient.grpcHandler.clientShellSystem.Get(grpcClient.Id, shellRequest.Id); err != nil {
		return
	}

	if clientShell == nil {
		err = ErrClientShellNotFound
		return
	}

	clientShell, err = grpcClient.grpcHandler.clientShellSystem.Save(grpcClient.Id, shellRequest.Id, func(clientShell *ClientShell) (rClientShell *ClientShell, err error) {
		rClientShell = clientShell.WithStatus(shellRequest.Status).WithOutput(clientShell.Output + string(shellRequest.Data))
		return
	})

	if err != nil {
		return
	}
	return
}
