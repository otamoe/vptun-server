package server

import (
	pb "github.com/otamoe/vptun-pb"
)

func (grpcClient *GrpcClient) OnPing(pingRequest *pb.PingRequest) (err error) {
	if !pingRequest.Pong {
		err = grpcClient.Response(&pb.StreamResponse{
			Ping: &pb.PingResponse{Pong: true},
		})
	}
	return
}

// SessionKey s
// sessionID := grpcClientLayerParsed.Type.String() + ":" + grpcClientLayerParsed.SourceIP.String() + ":" + strconv.FormatUint(uint64(grpcClientLayerParsed.SourcePort), 10)
