package server

import (
	pb "github.com/otamoe/vptun-pb"
	"go.uber.org/zap"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (grpcClient *GrpcClient) OnTun(tunRequest *pb.TunRequest) (err error) {
	if len(tunRequest.Data) == 0 {
		return
	}
	// 长度不正确
	if len(tunRequest.Data) > 4096 {
		err = grpc.Errorf(codes.InvalidArgument, "Unknown address version")
		return
	}

	if grpcClient.grpcHandler.subnetv6 {
		// ipv6
		if tunRequest.Data[0]>>4 != ipv6.Version {
			err = grpc.Errorf(codes.InvalidArgument, "Unknown address version")
			return
		}
	} else {
		// ipv4
		if tunRequest.Data[0]>>4 != ipv4.Version {
			err = grpc.Errorf(codes.InvalidArgument, "Unknown address version")
			return
		}
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

	if err != nil {
		return
	}

	// 是自己直接返回给自己
	if !grpcClient.IRouteAddress.Equal(grpcClientLayerParsed.DestinationIP) {
		err = grpcClient.Response(&pb.StreamResponse{Tun: &pb.TunResponse{Data: tunRequest.Data}})
		return
	}

	// 源 ip 不匹配
	if !grpcClient.IRouteAddress.Equal(grpcClientLayerParsed.SourceIP) {
		err = grpc.Errorf(codes.InvalidArgument, "Source IP address does not match")
		return
	}

	// 不在子网内
	if grpcClient.grpcHandler.subnet.Contains(grpcClientLayerParsed.DestinationIP) {
		err = grpc.Errorf(codes.InvalidArgument, "Destination IP address does not match")
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

	// 路由不通过
	if !grpcClient.routesOK(routes, grpcClientLayerParsed, dstGrpcClient) {
		return
	}

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

// parser := gopacket.NewDecodingLayerParser(layers.LayerTypeIPv4, &ip4, &tcp, &udp, &icmp)
// // decoded := []gopacket.LayerType{}
// // parser.DecodeLayers(b[4:], &decoded)

// // fmt.Println(ip4.SrcIP)
// // fmt.Println(ip4.DstIP)
// // fmt.Println(tcp.SrcPort)
// // fmt.Println(tcp.DstPort)
// // fmt.Println(udp.SrcPort)
// // fmt.Println(udp.DstPort)
// for _, layerType := range decoded {
// 	switch layerType {
// 	case layers.LayerTypeIPv4:
// 		fmt.Println("    IP4 ", ip4.SrcIP, ip4.DstIP)
// 	case layers.LayerTypeICMPv4:
// 		fmt.Println("    ICMP ")
// 	case layers.LayerTypeUDP:
// 		fmt.Println("    UDP ", udp.SrcPort, udp.DstPort)
// 	case layers.LayerTypeTCP:
// 		fmt.Println("    TCP ", tcp.SrcPort, tcp.DstPort)
// 	}
// }
// fmt.Println("   ")
// fmt.Println("   ")
