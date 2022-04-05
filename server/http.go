package server

import (
	"crypto/tls"
	"os"
	"path"

	libhttp "github.com/otamoe/go-library/http"
	libviper "github.com/otamoe/go-library/viper"
	"github.com/spf13/viper"
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configDir := path.Join(userHomeDir, ".vptun")
	viper.AddConfigPath("$HOME/.vptun/server/http")
	libviper.SetDefault("http.listenAddress", ":8443", "HTTP listen address")
	libviper.SetDefault("http.username", "admin", "Admin username")
	libviper.SetDefault("http.password", "", "Admin password")
	libviper.SetDefault("http.tlsCA", path.Join(configDir, "server/http/ca.crt"), "Certification authority")
	libviper.SetDefault("http.tlsCrt", path.Join(configDir, "server/http/server.crt"), "Credentials")
	libviper.SetDefault("http.tlsKey", path.Join(configDir, "server/http/server.key"), "Private key")
	libviper.SetDefault("http.public", path.Join(configDir, "public"), "public file directory")
}

func HttpTLSOption() (out libhttp.OutOption, err error) {
	var tlsConfig *tls.Config
	if tlsConfig, err = ParseTLS(viper.GetString("http.tlsCA"), viper.GetString("http.tlsCrt"), viper.GetString("http.tlsKey")); err != nil {
		return
	}
	tlsConfig.MinVersion = tls.VersionTLS12
	out.Option = func(server *libhttp.Server) error {
		server.TLSConfig = tlsConfig
		return nil
	}
	return
}

func HttpListenAddressOption() (out libhttp.OutOption) {
	listenAddress := viper.GetString("http.listenAddress")
	if listenAddress != "" {
		out = libhttp.WithListenAddress(listenAddress)()
	} else {
		out.Option = func(server *libhttp.Server) error {
			return nil
		}
	}
	return
}

func HttpInvoke(_ *libhttp.Server) (err error) {
	return
}
