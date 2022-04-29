package server

import (
	"errors"

	libbadger "github.com/otamoe/go-library/badger"
	libgrpc "github.com/otamoe/go-library/grpc"
	libhttp "github.com/otamoe/go-library/http"
	liblogger "github.com/otamoe/go-library/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var logger = liblogger.Get("server")

var ErrServerCA = errors.New("Server CA")
var ErrServerCertificate = errors.New("Server certificate")

func New() fx.Option {
	return fx.Options(
		libgrpc.New(),
		fx.Provide(GrpcEnforcementPolicyServerOption),
		fx.Provide(GrpcKeepaliveParamsServerOption),
		fx.Provide(GrpcMaxConcurrentStreamsServerOption),
		fx.Provide(GrpcConnectionTimeoutServerOption),
		fx.Provide(GrpcListenAddressOption),
		fx.Provide(GrpcTLSServerOption),
		fx.Provide(GrpcHandlerRegister),
		fx.Invoke(GrpcInvoke),

		libhttp.New(),
		fx.Provide(libhttp.WithErrorLog(zap.NewStdLog(liblogger.Get("http")))),
		fx.Provide(HttpListenAddressOption),
		fx.Provide(HttpTLSOption),
		fx.Provide(HttpHandlerRegister),
		fx.Invoke(HttpInvoke),

		libbadger.New(),
		fx.Provide(BadgerIndexDirOption),
		fx.Provide(BadgerValueDirOption),
		fx.Provide(BadgerBlockCacheSizeOption),
		fx.Provide(BadgerIndexCacheSizeOption),
		fx.Provide(BadgerMemTableSizeOption),

		fx.Provide(NewClientSystem),
		fx.Provide(NewRouteSystem),
		fx.Provide(NewClientShellSystem),
	)
}
