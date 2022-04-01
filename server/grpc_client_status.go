package server

import (
	pb "github.com/otamoe/vptun-pb"
)

func (grpcClient *GrpcClient) OnStatus(statusRequest *pb.StatusRequest) (err error) {
	if statusRequest.Status == nil {
		return
	}
	_, err = grpcClient.grpcHandler.clientSystem.Save(grpcClient.Id, func(client *Client) (rClient *Client, err error) {
		rClient = client.WithStatus(statusRequest.Status)
		return
	})

	if err != nil {
		return
	}
	err = grpcClient.Response(&pb.StreamResponse{
		Status: &pb.StatusResponse{},
	})
	return
}
