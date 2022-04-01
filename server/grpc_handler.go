package server

import (
	"context"
	"strings"
	"sync"
	"time"

	libgrpc "github.com/otamoe/go-library/grpc"
	liblogger "github.com/otamoe/go-library/logger"
	pb "github.com/otamoe/vptun-pb"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type (
	GrpcHandler struct {
		routeSystem           *RouteSystem
		clientSystem          *ClientSystem
		clientShellSystem     *ClientShellSystem
		mux                   sync.RWMutex
		ctx                   context.Context
		cancel                context.CancelFunc
		clientsByRouteAddress map[string]*GrpcClient
		clients               map[string]*GrpcClient
		routes                Routes
		logger                *zap.Logger
	}
)

func GrpcHandlerRegister(ctx context.Context, routeSystem *RouteSystem, clientSystem *ClientSystem, clientShellSystem *ClientShellSystem, lc fx.Lifecycle) (grpcHandler *GrpcHandler, out libgrpc.OutServer, err error) {
	ctx, cancel := context.WithCancel(ctx)
	grpcHandler = &GrpcHandler{
		logger:                liblogger.Get("grpc"),
		routeSystem:           routeSystem,
		clientSystem:          clientSystem,
		clientShellSystem:     clientShellSystem,
		ctx:                   ctx,
		cancel:                cancel,
		routes:                routeSystem.All(),
		clientsByRouteAddress: map[string]*GrpcClient{},
		clients:               map[string]*GrpcClient{},
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// 5 秒更新 一次路由表
			go func() {
				t := time.NewTicker(time.Second * 5)
				defer t.Stop()
				select {
				case <-t.C:
					grpcHandler.mux.Lock()
					grpcHandler.routes = grpcHandler.routeSystem.All()
					grpcHandler.mux.Unlock()
				case <-grpcHandler.ctx.Done():
					return
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancel()
			return nil
		},
	})
	out.Server = func(s *grpc.Server) (err error) {
		pb.RegisterHandlerServer(s, grpcHandler)
		return nil
	}
	return
}

func (grpcHandler *GrpcHandler) readPeer(ctx context.Context) (peerFrom *peer.Peer, md metadata.MD, err error) {
	var ok bool
	var tlsInfo credentials.TLSInfo

	// 读取 peer 信息
	peerFrom, ok = peer.FromContext(ctx)

	if !ok {
		err = grpc.Errorf(codes.PermissionDenied, "peer not found")
		return
	}
	if peerFrom.Addr == nil {
		err = grpc.Errorf(codes.PermissionDenied, "address not found")
		return
	}

	if peerFrom.AuthInfo == nil || peerFrom.AuthInfo.AuthType() != "tls" {
		err = grpc.Errorf(codes.PermissionDenied, "tls not found")
		return
	}

	// 读取 metadata 信息
	md, ok = metadata.FromIncomingContext(ctx)
	if !ok {
		err = grpc.Errorf(codes.PermissionDenied, "metadata not found")
		return
	}

	// 读取 证书
	tlsInfo, ok = peerFrom.AuthInfo.(credentials.TLSInfo)
	if !ok {
		err = grpc.Errorf(codes.PermissionDenied, "tls not found")
		return
	}

	// 检查 证书
	if len(tlsInfo.State.VerifiedChains) == 0 || len(tlsInfo.State.VerifiedChains[0]) == 0 {
		err = grpc.Errorf(codes.PermissionDenied, "certificate server name does not match")
		return
	}

	// 证书不是客户端
	if tlsInfo.State.VerifiedChains[0][0].Subject.CommonName != "client" {
		err = grpc.Errorf(codes.PermissionDenied, "certificate server name does not match")
		return
	}
	return
}

func (grpcHandler *GrpcHandler) readPeerClient(ctx context.Context) (client *Client, peerFrom *peer.Peer, md metadata.MD, err error) {
	if peerFrom, md, err = grpcHandler.readPeer(ctx); err != nil {
		return
	}

	client = grpcHandler.clientSystem.Get(strings.Join(md.Get("client-id"), ","))

	// 客户端找不到
	if client == nil {
		err = grpc.Errorf(codes.PermissionDenied, "client not found")
		return
	}

	// 密码不正确
	if client.Key == "" || client.Key != strings.Join(md.Get("client-key"), ",") {
		err = grpc.Errorf(codes.PermissionDenied, "client key does not match")
		return
	}

	// 状态 不可用
	if client.State != pb.State_AVAILABLE {
		err = grpc.Errorf(codes.PermissionDenied, "client is unavailable")
		return
	}

	// 已过期
	if time.Now().Unix() > client.ExpiredAt {
		err = grpc.Errorf(codes.PermissionDenied, "client has expired")
		return
	}
	return
}

func (grpcHandler *GrpcHandler) GetClientByRouteAddress(routeAddress string) (client *GrpcClient) {
	grpcHandler.mux.RLock()
	defer grpcHandler.mux.RUnlock()
	client, _ = grpcHandler.clientsByRouteAddress[routeAddress]
	return
}

func (grpcHandler *GrpcHandler) GetClient(id string) (client *GrpcClient) {
	grpcHandler.mux.RLock()
	defer grpcHandler.mux.RUnlock()
	client, _ = grpcHandler.clients[id]
	return
}
