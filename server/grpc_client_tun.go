package server

import (
	"time"

	pb "github.com/otamoe/vptun-pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (grpcClient *GrpcClient) OnTun(tunRequest *pb.TunRequest) (err error) {
	if len(tunRequest.Data) == 0 {
		return
	}
	// 长度不正确
	if len(tunRequest.Data) > 4096 {
		err = grpc.Errorf(codes.InvalidArgument, "Packet too long")
		return
	}

	// 解析包
	var grpcClientLayerParsed *GrpcClientLayerParsed
	grpcClientLayerParsed, err = grpcClient.layerParse(tunRequest.Data)

	fields := []zap.Field{}
	if grpcClientLayerParsed != nil {
		fields = append(
			fields,
			zap.String("type", grpcClientLayerParsed.Type.String()),
			zap.String("sourceIP", grpcClientLayerParsed.SourceIP.String()),
			zap.String("destinationIP", grpcClientLayerParsed.DestinationIP.String()),
			zap.Uint32("sourcePort", grpcClientLayerParsed.SourcePort),
			zap.Uint32("destinationPort", grpcClientLayerParsed.DestinationPort),
		)
	}
	grpcClient.logger("tun", true, err, fields...)

	// 解析包错误
	if err != nil {
		err = nil
		return
	}

	// 是自己直接返回给自己
	if grpcClient.IRouteAddress.Equal(grpcClientLayerParsed.DestinationIP) {
		err = grpcClient.Response(&pb.StreamResponse{Tun: &pb.TunResponse{Data: tunRequest.Data}})
		return
	}

	// 源 ip 不匹配
	if !grpcClient.IRouteAddress.Equal(grpcClientLayerParsed.SourceIP) {
		return
	}

	// 不在子网内
	if !grpcClient.grpcHandler.subnet.Contains(grpcClientLayerParsed.DestinationIP) {
		return
	}

	// 路由
	grpcClient.grpcHandler.mux.RLock()
	dstGrpcClient, _ := grpcClient.grpcHandler.clientsByRouteAddress[grpcClientLayerParsed.DestinationIP.String()]
	routes := grpcClient.grpcHandler.routes
	grpcClient.grpcHandler.mux.RUnlock()

	//  目标不存在
	if dstGrpcClient == nil {
		return
	}

	sessionKey := SessionKey{
		Type:            grpcClientLayerParsed.Type,
		SourcePort:      grpcClientLayerParsed.SourcePort,
		DestinationIP:   grpcClientLayerParsed.DestinationIP.String(),
		DestinationPort: grpcClientLayerParsed.DestinationPort,
	}

	// 是否路由通过
	var ok bool

	// 查找 目标方 session
	{
		dstGrpcClient.mux.Lock()
		dstSessionKey := SessionKey{
			Type:            grpcClientLayerParsed.Type,
			SourcePort:      grpcClientLayerParsed.DestinationPort,
			DestinationIP:   grpcClientLayerParsed.SourceIP.String(),
			DestinationPort: grpcClientLayerParsed.SourcePort,
		}
		if _, ok = dstGrpcClient.sessions[dstSessionKey]; ok {
			dstGrpcClient.sessions[dstSessionKey] = time.Now()
		}
		dstGrpcClient.mux.Unlock()
	}

	// 路由不通过
	if !ok && !grpcClient.routesOK(routes, grpcClientLayerParsed, dstGrpcClient) {
		return
	}

	// 储存 自己的 sessions
	grpcClient.mux.Lock()
	grpcClient.sessions[sessionKey] = time.Now()
	grpcClient.mux.Unlock()

	// 发送给对方
	err = dstGrpcClient.Response(&pb.StreamResponse{Tun: &pb.TunResponse{Data: tunRequest.Data}})

	return
}

func (grpcClient *GrpcClient) routesOK(routes Routes, grpcClientLayerParsed *GrpcClientLayerParsed, dstGrpcClient *GrpcClient) (ok bool) {
	for _, route := range routes {
		// 不是激活的
		if route.State != pb.State_AVAILABLE {
			continue
		}

		// 协议 匹配
		if route.Type != pb.Route_NONE && grpcClientLayerParsed.Type != route.Type {
			continue
		}

		// 源 ip
		if !route.iSourceIP.Contains(grpcClientLayerParsed.SourceIP) {
			continue
		}

		// 目标 ip
		if !route.iDestinationIP.Contains(grpcClientLayerParsed.DestinationIP) {
			continue
		}

		// 源 port
		if route.SourcePort != 0 && grpcClientLayerParsed.SourcePort != route.SourcePort {
			continue
		}

		// 目标 port
		if route.DestinationPort != 0 && grpcClientLayerParsed.DestinationPort != route.DestinationPort {
			continue
		}

		// 不是 拒绝 就接受
		ok = route.Action != pb.Route_REJECT
		break
	}

	// 储存 session
	// sessionID := grpcClientLayerParsed.Type.String() + ":" + grpcClientLayerParsed.SourceIP.String() + ":" + strconv.FormatUint(uint64(grpcClientLayerParsed.SourcePort), 10)
	// dstGrpcClient.session
	return
}

// SessionKey s
// sessionID := grpcClientLayerParsed.Type.String() + ":" + grpcClientLayerParsed.SourceIP.String() + ":" + strconv.FormatUint(uint64(grpcClientLayerParsed.SourcePort), 10)
