package server

import (
	pb "github.com/otamoe/vptun-pb"
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

	if tunRequest.Data[0]>>4 == ipv6.Version {
		if err = grpcClient.tunRecvIPV6(tunRequest.Data); err != nil {
			return
		}
	} else {
		if err = grpcClient.tunRecvIPV4(tunRequest.Data); err != nil {
			return
		}
	}

	return
}

func (grpcClient *GrpcClient) tunRecvIPV6(data []byte) (err error) {
	// ipv4
	var header *ipv6.Header

	if header, err = ipv6.ParseHeader(data); err != nil {
		return
	}

	// src 地址 不匹配
	if !header.Src.Equal(grpcClient.IRouteAddress) {
		err = grpc.Errorf(codes.InvalidArgument, "Source address does not match")
		return
	}

	// 本地
	if header.Dst.Equal(grpcClient.IRouteAddress) {
		err = grpcClient.Response(&pb.StreamResponse{Tun: &pb.TunResponse{Data: data}})
		return
	}

	// 路由
	grpcClient.grpcHandler.mux.RLock()
	dstGrpcClient, _ := grpcClient.grpcHandler.clientsByRouteAddress[header.Dst.String()]
	routes := grpcClient.grpcHandler.routes
	grpcClient.grpcHandler.mux.RUnlock()

	if dstGrpcClient != nil {
		var ok bool
		for _, route := range routes {
			// 不是激活的
			if route.State != pb.State_AVAILABLE {
				continue
			}

			// 源 ip
			if !route.iSource.Contains(header.Src) {
				continue
			}

			// 目标 ip
			if !route.iDestination.Contains(header.Dst) {
				continue
			}

			// 不是 拒绝 就接受
			ok = route.Action != pb.Route_REJECT

			break
		}
		if ok {
			err = dstGrpcClient.Response(&pb.StreamResponse{Tun: &pb.TunResponse{Data: data}})
		}
		return
	}
	return
}
func (grpcClient *GrpcClient) tunRecvIPV4(data []byte) (err error) {
	// ipv4
	var header *ipv4.Header

	if header, err = ipv4.ParseHeader(data); err != nil {
		return
	}

	// src 地址 不匹配
	if !header.Src.Equal(grpcClient.IRouteAddress) {
		err = grpc.Errorf(codes.InvalidArgument, "Source address does not match")
		return
	}

	// 本地
	if header.Dst.Equal(grpcClient.IRouteAddress) {
		err = grpcClient.Response(&pb.StreamResponse{Tun: &pb.TunResponse{Data: data}})
		return
	}

	// 路由
	grpcClient.grpcHandler.mux.RLock()
	dstGrpcClient, _ := grpcClient.grpcHandler.clientsByRouteAddress[header.Dst.String()]
	routes := grpcClient.grpcHandler.routes
	grpcClient.grpcHandler.mux.RUnlock()

	if dstGrpcClient != nil {
		var ok bool
		for _, route := range routes {
			// 不是激活的
			if route.State != pb.State_AVAILABLE {
				continue
			}

			// 源 ip
			if !route.iSource.Contains(header.Src) {
				continue
			}

			// 目标 ip
			if !route.iDestination.Contains(header.Dst) {
				continue
			}

			// 不是 拒绝 就接受
			ok = route.Action != pb.Route_REJECT

			break
		}
		if ok {
			err = dstGrpcClient.Response(&pb.StreamResponse{Tun: &pb.TunResponse{Data: data}})
		}
		return
	}
	return
}
