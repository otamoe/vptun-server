package server

import (
	"crypto/tls"
	"os"
	"path"
	"time"

	libgrpc "github.com/otamoe/go-library/grpc"
	libviper "github.com/otamoe/go-library/viper"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configDir := path.Join(userHomeDir, ".vptun")

	libviper.SetDefault("grpc.listenAddress", ":9443", "Listen address")
	libviper.SetDefault("grpc.maxConnectionIdle", time.Minute*5, "Max connection idle")
	libviper.SetDefault("grpc.maxConnectionAge", time.Duration(0), "Max connection age")
	libviper.SetDefault("grpc.maxConnectionAgeGrace", time.Duration(0), "Max connection age grace")
	libviper.SetDefault("grpc.time", time.Minute*4, "Ping interval")
	libviper.SetDefault("grpc.timeout", time.Second*30, "Ping timeout")

	libviper.SetDefault("grpc.connectionTimeout", time.Second*30, "Connection timeout")

	libviper.SetDefault("grpc.maxConcurrentStreams", uint32(64), "Max concurrent streams")
	libviper.SetDefault("grpc.tlsCA", path.Join(configDir, "server/grpc/ca.crt"), "Certification authority")
	libviper.SetDefault("grpc.tlsCrt", path.Join(configDir, "server/grpc/server.crt"), "Credentials")
	libviper.SetDefault("grpc.tlsKey", path.Join(configDir, "server/grpc/server.key"), "Private key")
}

func GrpcKeepaliveParamsServerOption() (out libgrpc.OutServerOption) {
	out.Option = grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     viper.GetDuration("grpc.maxConnectionIdle"),
		MaxConnectionAge:      viper.GetDuration("grpc.maxConnectionAge"),
		MaxConnectionAgeGrace: viper.GetDuration("grpc.maxConnectionAgeGrace"),
		Time:                  viper.GetDuration("grpc.time"),
		Timeout:               viper.GetDuration("grpc.timeout"),
	})
	return
}

func GrpcMaxConcurrentStreamsServerOption() (out libgrpc.OutServerOption) {
	out.Option = grpc.MaxConcurrentStreams(viper.GetUint32("grpc.maxConcurrentStreams"))
	return
}

func GrpcConnectionTimeoutServerOption() (out libgrpc.OutServerOption) {
	out.Option = grpc.ConnectionTimeout(viper.GetDuration("grpc.connectionTimeout"))
	return
}

func GrpcTLSServerOption() (out libgrpc.OutServerOption, err error) {
	var tlsConfig *tls.Config
	if tlsConfig, err = ParseTLS(viper.GetString("grpc.tlsCA"), viper.GetString("grpc.tlsCrt"), viper.GetString("grpc.tlsKey")); err != nil {
		return
	}
	// tls 配置
	tlsConfig.ServerName = "server"
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	tlsConfig.ClientCAs = tlsConfig.RootCAs

	out.Option = grpc.Creds(credentials.NewTLS(tlsConfig))
	return
}

func GrpcInvoke(_ *grpc.Server) (err error) {
	return
}

func GrpcListenAddressOption() (out libgrpc.OutExtendedServerOption) {
	out.Option = func(extendedServerOptions *libgrpc.ExtendedServerOptions) (err error) {
		extendedServerOptions.ListenAddress = viper.GetString("grpc.listenAddress")
		return nil
	}
	return
}
