package server

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	pb "github.com/otamoe/vptun-pb"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func (grpcHandler *GrpcHandler) Stream(handlerStreamServer pb.Handler_StreamServer) (err error) {
	var peerFrom *peer.Peer
	var md metadata.MD
	var client *Client
	connectAt := time.Now()
	wLogger := func(name string, debug bool, err error, fields ...zap.Field) {
		if md != nil {
			fields = append(
				fields,
				zap.String("userAgent", strings.Join(md.Get("user-agent"), ",")),
				zap.String("clientHostname", strings.Join(md.Get("client-hostname"), ",")),
				zap.String("clientShell", strings.Join(md.Get("client-shell"), ",")),
				zap.Duration("duration", time.Now().Sub(connectAt)),
			)
		}
		if peerFrom != nil {
			fields = append(
				fields,
				zap.String("remoteAddress", peerFrom.Addr.String()),
			)
		}
		if client != nil {
			fields = append(
				fields,
				zap.String("clientID", client.Id),
				zap.String("clientRouteAddress", client.RouteAddress),
			)
		}
		if err != nil {
			fields = append(fields, zap.Error(err))
			if IsErrClose(err) {
				grpcHandler.logger.Info(name, fields...)
			} else if IsErrClient(err) {
				grpcHandler.logger.Warn(name, fields...)
			} else {
				fields = append(fields, zap.Stack("stack"))
				grpcHandler.logger.Error(name, fields...)
			}
		} else {
			if debug {
				grpcHandler.logger.Debug(name, fields...)
			} else {
				grpcHandler.logger.Info(name, fields...)
			}
		}
	}

	defer func() {
		wLogger("stream-end", false, err)
	}()

	client, peerFrom, md, err = grpcHandler.readPeerClient(handlerStreamServer.Context())
	wLogger("stream-run", false, err)
	if err != nil {
		return
	}

	// 发送 header
	if err = handlerStreamServer.SendHeader(metadata.New(map[string]string{
		"client-expired-at":    strconv.FormatInt(client.ExpiredAt, 10),
		"client-route-address": client.RouteAddress,
		"client-route-subnet":  grpcHandler.subnet.String(),
	})); err != nil {
		return
	}

	// 检查 连接池 存在 并且关闭他
	grpcClient := grpcHandler.GetClient(client.Id)
	{
		if grpcClient != nil {
			grpcClient.Close(grpc.Errorf(codes.Canceled, "New client connected"))
		}

		grpcClient = &GrpcClient{
			Id:                  client.Id,
			RouteAddress:        client.RouteAddress,
			IRouteAddress:       client.iRouteAddress,
			MD:                  md,
			Peer:                peerFrom,
			ConnectAt:           connectAt,
			handlerStreamServer: handlerStreamServer,
			grpcHandler:         grpcHandler,
			logger:              wLogger,
			sessions:            make(map[SessionKey]time.Time),
			response:            make(chan *pb.StreamResponse, 4),
			close:               atomic.NewBool(false),
			closed:              make(chan struct{}),
			shells:              map[string]*GrpcClientShell{},
		}
		if grpcHandler.subnetv6 {
			grpcClient.layerParser = gopacket.NewDecodingLayerParser(layers.LayerTypeIPv6, &grpcClient.layerIPv6, &grpcClient.layerTCP, &grpcClient.layerUDP, &grpcClient.layerICMPv6)
		} else {
			grpcClient.layerParser = gopacket.NewDecodingLayerParser(layers.LayerTypeIPv4, &grpcClient.layerIPv4, &grpcClient.layerTCP, &grpcClient.layerUDP, &grpcClient.layerICMPv4)
		}
		grpcClient.ctx, grpcClient.cancel = context.WithCancel(handlerStreamServer.Context())
	}

	// 关闭
	defer func() {
		if rerr := grpcClient.Close(err); rerr != nil && err == nil {
			err = rerr
		}
		if err == nil && grpcClient.Err() != os.ErrClosed {
			err = grpcClient.Err()
		}
	}()

	// 已取消
	defer close(grpcClient.closed)

	// 写入连接池
	{
		grpcHandler.mux.Lock()
		if _, ok := grpcHandler.clients[client.Id]; ok {
			grpcHandler.mux.Unlock()
			// 已存在
			err = grpc.Errorf(codes.AlreadyExists, "Client already connected")
			return
		}

		grpcHandler.clientsByRouteAddress[client.RouteAddress] = grpcClient
		grpcHandler.clients[client.Id] = grpcClient
		grpcHandler.mux.Unlock()
	}

	// 读取数据
	go grpcClient.runRequest()

	// 写入数据
	go grpcClient.runResponse()

	// 运行在线
	go grpcClient.runOnline()

	// 运行 Session
	go grpcClient.runSession()

	// 结束
	<-grpcClient.ctx.Done()
	return
}
