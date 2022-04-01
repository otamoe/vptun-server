package server

import (
	"context"
	"net"
	"os"
	"strings"
	"sync"
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

type (
	GrpcClient struct {
		Id            string
		RouteAddress  string
		IRouteAddress net.IP

		MD                  metadata.MD
		Peer                *peer.Peer
		ConnectAt           time.Time
		grpcHandler         *GrpcHandler
		handlerStreamServer pb.Handler_StreamServer
		logger              func(name string, debug bool, err error, fields ...zap.Field)
		mux                 sync.RWMutex
		err                 error
		ctx                 context.Context
		cancel              context.CancelFunc
		close               *atomic.Bool
		closed              chan struct{}

		response chan *pb.StreamResponse
		shells   map[string]*GrpcClientShell

		layerIPv4   layers.IPv4
		layerIPv6   layers.IPv6
		layerICMPv4 layers.ICMPv4
		layerICMPv6 layers.ICMPv6
		layerTCP    layers.TCP
		layerUDP    layers.UDP
		layerParser *gopacket.DecodingLayerParser
	}
)

// 错误
func (grpcClient *GrpcClient) Err() (err error) {
	select {
	case <-grpcClient.closed:
		err = grpcClient.err
		return
	}
}

func (grpcClient *GrpcClient) runResponse() (err error) {
	defer func() {
		go grpcClient.Close(err)
	}()

	var response *pb.StreamResponse
	for {
		select {
		case <-grpcClient.ctx.Done():
			return
		case response = <-grpcClient.response:
			if response == nil {
				return
			}
			err = grpcClient.handlerStreamServer.Send(response)
			fields := []zap.Field{}
			if response.Tun != nil {
				fields = append(
					fields,
					zap.Int("requestTunData", len(response.Tun.Data)),
				)
			}
			if response.Shell != nil {
				fields = append(
					fields,
					zap.String("responseShellId", response.Shell.Id),
					zap.String("responseShellData", string(response.Shell.Data)),
					zap.Uint32("responseShellTimeout", response.Shell.Timeout),
				)
			}
			if response.Status != nil {
				fields = append(
					fields,
					zap.Bool("responseStatus", true),
				)
			}
			grpcClient.logger("stream-send", response.Shell == nil, err, fields...)
		}
	}
}

func (grpcClient *GrpcClient) runRequest() (err error) {
	defer func() {
		go grpcClient.Close(err)
	}()

	var request *pb.StreamRequest
	for {
		request, err = grpcClient.handlerStreamServer.Recv()
		fields := []zap.Field{}
		if request != nil {
			if request.Tun != nil {
				fields = append(
					fields,
					zap.Int("requestTunData", len(request.Tun.Data)),
				)
				if err == nil {
					err = grpcClient.OnTun(request.Tun)
				}
			}
			if request.Shell != nil {
				fields = append(
					fields,
					zap.String("requestShellId", request.Shell.Id),
					zap.String("requestShellData", string(request.Shell.Data)),
					zap.Int32("requestShellStatus", request.Shell.Status),
				)
				if err == nil {
					err = grpcClient.OnShell(request.Shell)
				}
			}
			if request.Status != nil {
				fields = append(
					fields,
					zap.Bool("requestStatus", true),
				)
				if request.Status.Status != nil && request.Status.Status.Load != nil {
					fields = append(
						fields,
						zap.Float64("requestStatusLoad1", request.Status.Status.Load.Load1),
						zap.Float64("requestStatusLoad5", request.Status.Status.Load.Load5),
						zap.Float64("requestStatusLoad15", request.Status.Status.Load.Load15),
					)
				}
				if err == nil {
					err = grpcClient.OnStatus(request.Status)
				}
			}
		}

		// log
		grpcClient.logger("stream-recv", request != nil && request.Shell == nil, err, fields...)
		if err != nil {
			return
		}
	}
}

func (grpcClient *GrpcClient) runOnline() (err error) {
	defer func() {
		go grpcClient.Close(err)
	}()

	hostname := strings.Join(grpcClient.MD.Get("client-hostname"), ",")
	userAgent := strings.Join(grpcClient.MD.Get("user-agent"), ",")
	shell := strings.Join(grpcClient.MD.Get("client-shell"), ",") == "1" || strings.Join(grpcClient.MD.Get("client-shell"), ",") == "true"
	var newClient *Client
	var i = 0

	t := time.NewTicker(time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-grpcClient.ctx.Done():
			return
		case <-t.C:
			t.Reset(time.Second * 30)
			// 设置 在线 hostname user-agent 在线 连接地址
			newClient, err = grpcClient.grpcHandler.clientSystem.Save(grpcClient.Id, func(client *Client) (rClient *Client, err error) {
				rClient = client.
					WithHostname(hostname).
					WithUserAgent(userAgent).
					WithConnectAddress(grpcClient.Peer.Addr.String()).
					WithShell(shell).
					WithOnline(true)
				if i == 0 {
					rClient = rClient.WithConnectAt(grpcClient.ConnectAt)
				}
				return
			})
			// 更新 错误
			if err != nil {
				return
			}

			// 路由地址 已修改
			if newClient.Client.RouteAddress != grpcClient.RouteAddress {
				err = grpc.Errorf(codes.Canceled, "Routing address has been updated")
				return
			}
			i++
		}
	}
}

func (grpcClient *GrpcClient) Response(response *pb.StreamResponse) (err error) {
	select {
	case <-grpcClient.ctx.Done():
		// 上下文取消
		err = grpcClient.ctx.Err()
	case <-grpcClient.closed:
		// 已关闭
		err = grpcClient.err
	case grpcClient.response <- response:
		// 写入
		if response.Shell != nil {
			if response.Shell.Timeout != 0 && len(response.Shell.Data) != 0 {
				// 创建
				grpcClient.mux.Lock()
				grpcClientShell := &GrpcClientShell{
					Id:         response.Shell.Id,
					grpcClient: grpcClient,
					close:      atomic.NewBool(false),
				}
				grpcClient.shells[response.Shell.Id] = grpcClientShell
				grpcClient.mux.Unlock()

				fields := []zap.Field{
					zap.String("clientId", grpcClient.Id),
					zap.String("shellId", grpcClientShell.Id),
					zap.String("data", string(response.Shell.Data)),
					zap.Uint32("timeout", response.Shell.Timeout),
				}
				grpcClient.logger("shell-open", false, nil, fields...)
			} else {
				// 取消
				grpcClientShell := grpcClient.Shell(response.Shell.Id)
				if grpcClientShell != nil {
					err = grpcClientShell.Close(1, context.Canceled)
					if err == os.ErrClosed {
						err = nil
					}
				}
			}
		}
	}
	return
}

func (grpcClient *GrpcClient) Shell(id string) (grpcClientShell *GrpcClientShell) {
	grpcClient.mux.RLock()
	grpcClientShell, _ = grpcClient.shells[id]
	grpcClient.mux.RUnlock()
	return
}

func (grpcClient *GrpcClient) Close(werr error) (rerr error) {
	if !grpcClient.close.CAS(false, true) {
		rerr = grpcClient.Err()
		return
	}

	if werr == nil {
		werr = os.ErrClosed
	}

	// 删除连接池
	grpcClient.grpcHandler.mux.Lock()
	if c, ok := grpcClient.grpcHandler.clients[grpcClient.Id]; ok && c == grpcClient {
		delete(grpcClient.grpcHandler.clients, grpcClient.Id)
	}
	if c, ok := grpcClient.grpcHandler.clientsByRouteAddress[grpcClient.RouteAddress]; ok && c == grpcClient {
		delete(grpcClient.grpcHandler.clientsByRouteAddress, grpcClient.RouteAddress)
	}
	grpcClient.grpcHandler.mux.Unlock()

	// 设置成离线
	_, rerr = grpcClient.grpcHandler.clientSystem.Save(grpcClient.Id, func(client *Client) (rClient *Client, err error) {
		rClient = client.WithOnline(false)
		return
	})

	grpcClientShells := []*GrpcClientShell{}
	grpcClient.mux.RLock()
	for _, grpcClientShell := range grpcClient.shells {
		grpcClientShells = append(grpcClientShells, grpcClientShell)
	}
	grpcClient.mux.RUnlock()
	for _, grpcClientShell := range grpcClientShells {
		if shellerr := grpcClientShell.Close(1, context.Canceled); shellerr != nil && shellerr != os.ErrClosed && rerr == nil {
			rerr = shellerr
		}
	}

	grpcClient.err = werr

	// 申请关闭
	grpcClient.cancel()

	// 连接已关闭
	<-grpcClient.closed

	return
}
