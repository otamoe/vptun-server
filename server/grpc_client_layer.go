package server

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	pb "github.com/otamoe/vptun-pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type (
	GrpcClientLayerParsed struct {
		SourceIP        net.IP
		DestinationIP   net.IP
		SourcePort      uint32
		DestinationPort uint32
		Type            pb.Route_Type
	}
)

func (grpcClient *GrpcClient) layerParse(data []byte) (grpcClientLayerParsed *GrpcClientLayerParsed, err error) {
	decoded := []gopacket.LayerType{}
	grpcClient.layerParser.DecodeLayers(data, &decoded)
	grpcClientLayerParsed = &GrpcClientLayerParsed{}
	for _, layerType := range decoded {
		switch layerType {
		case layers.LayerTypeIPv4:
			grpcClientLayerParsed.SourceIP = grpcClient.layerIPv4.SrcIP
			grpcClientLayerParsed.DestinationIP = grpcClient.layerIPv4.DstIP
		case layers.LayerTypeIPv6:
			grpcClientLayerParsed.SourceIP = grpcClient.layerIPv6.SrcIP
			grpcClientLayerParsed.DestinationIP = grpcClient.layerIPv6.DstIP
		case layers.LayerTypeICMPv4:
			grpcClientLayerParsed.Type = pb.Route_ICMP
		case layers.LayerTypeICMPv6:
			grpcClientLayerParsed.Type = pb.Route_ICMP
		case layers.LayerTypeUDP:
			grpcClientLayerParsed.Type = pb.Route_UDP
			grpcClientLayerParsed.SourcePort = uint32(uint16(grpcClient.layerUDP.SrcPort))
			grpcClientLayerParsed.DestinationPort = uint32(uint16(grpcClient.layerUDP.DstPort))
		case layers.LayerTypeTCP:
			grpcClientLayerParsed.Type = pb.Route_TCP
			grpcClientLayerParsed.SourcePort = uint32(uint16(grpcClient.layerTCP.SrcPort))
			grpcClientLayerParsed.DestinationPort = uint32(uint16(grpcClient.layerTCP.DstPort))
		}
	}
	if grpcClientLayerParsed.Type == pb.Route_NONE || grpcClientLayerParsed.SourceIP.IsUnspecified() || grpcClientLayerParsed.DestinationIP.IsUnspecified() {
		grpcClientLayerParsed = nil
		err = grpc.Errorf(codes.InvalidArgument, "Unknown address version")
	}
	return
}
