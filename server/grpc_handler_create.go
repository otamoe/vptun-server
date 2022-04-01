package server

import (
	"context"
	"strings"
	"time"

	libutils "github.com/otamoe/go-library/utils"
	pb "github.com/otamoe/vptun-pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func (grpcHandler *GrpcHandler) Create(ctx context.Context, createRequest *pb.CreateRequest) (createResponse *pb.CreateResponse, err error) {
	var md metadata.MD
	var peerFrom *peer.Peer
	var client *Client
	now := time.Now()
	wLogger := func(name string, fields ...zap.Field) {
		if md != nil {
			fields = append(
				fields,
				zap.String("userAgent", strings.Join(md.Get("user-agent"), ",")),
				zap.String("clientHostname", strings.Join(md.Get("client-hostname"), ",")),
				zap.Duration("duration", time.Now().Sub(now)),
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
				zap.String("clientId", client.Id),
				zap.String("clientRouteAddress", client.RouteAddress),
			)
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			if IsErrClose(err) {
				grpcHandler.logger.Info(name, fields...)
			} else if IsErrClose(err) {
				fields = append(fields, zap.Stack("stack"))
				grpcHandler.logger.Warn(name, fields...)
			} else {
				fields = append(fields, zap.Stack("stack"))
				grpcHandler.logger.Error(name, fields...)
			}
		} else {
			grpcHandler.logger.Info(name, fields...)
		}
	}

	defer wLogger("create-end")

	if peerFrom, md, err = grpcHandler.readPeer(ctx); err != nil {
		return
	}

	wLogger("create-run")

	client, err = grpcHandler.clientSystem.Save("", func(client *Client) (rClient *Client, err error) {
		routeAddress := grpcHandler.clientSystem.NewRouteAddress(false)
		if len(routeAddress) == 0 {
			err = grpc.Errorf(codes.PermissionDenied, "Client creation is not allowed")
			return
		}
		rClient = client.
			WithKey(string(libutils.RandByte(16, libutils.RandAlphaNumber))).
			WithHostname(strings.Join(md.Get("client-hostname"), ",")).
			WithUserAgent(strings.Join(md.Get("user-agent"), ",")).
			WithRouteAddress(routeAddress).
			WithState(pb.State_AVAILABLE).
			WithExpiredAt(time.Date(9001, time.January, 1, 0, 0, 0, 0, time.UTC))
		return
	})
	if err != nil {
		return
	}

	createResponse = &pb.CreateResponse{
		Client: client.Client,
	}

	return
}
