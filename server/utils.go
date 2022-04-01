package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	libutils "github.com/otamoe/go-library/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ParseTLS(ca string, crt string, key string) (tlsConfig *tls.Config, err error) {
	// tls 配置
	tlsConfig = &tls.Config{
		MinVersion: tls.VersionTLS13, // 最低版本 1.3
	}

	// server ca
	if ca != "" {
		certPool := x509.NewCertPool()
		var b []byte
		if b, err = ioutil.ReadFile(ca); err != nil {
			b = []byte(ca)
		}
		if !certPool.AppendCertsFromPEM(b) {
			err = ErrServerCA
			return
		}
		err = nil

		// root ca
		tlsConfig.RootCAs = certPool
	}

	// server crt, key
	var certificate tls.Certificate
	if certificate, err = tls.LoadX509KeyPair(crt, key); err != nil {
		if certificate, err = tls.X509KeyPair([]byte(crt), []byte(key)); err != nil {
			err = ErrServerCertificate
			return
		}
	}
	tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)

	return
}

func NewID(now time.Time) string {
	t := libutils.MarshalTime(now.UTC())
	buf := make([]byte, base64.RawURLEncoding.EncodedLen(len(t)))
	base64.RawURLEncoding.Encode(buf, t)
	buf = append(buf, libutils.RandByte(16, libutils.RandAlphaNumber)...)
	return string(buf)
}

func IsID(id string) bool {
	blen := base64.RawURLEncoding.EncodedLen(12)
	if len(id) != blen+16 {
		return false
	}

	for _, v := range id {
		if (v < '0' || v > '9') && (v < 'a' || v > 'z') && (v < 'A' || v > 'Z') && v != '_' && v != '-' {
			return false
		}
	}
	return true
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func IsErrClose(err error) bool {
	if err == nil {
		return false
	}
	if err == io.EOF || err.Error() == io.EOF.Error() {
		return true
	}
	if err == os.ErrClosed || err.Error() == os.ErrClosed.Error() {
		return true
	}
	if errors.Is(err, context.Canceled) {
		return true
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	if errors.Is(err, grpc.ErrClientConnClosing) {
		return true
	}

	if code := status.Code(err); code == codes.Canceled || code == codes.DeadlineExceeded {
		return true
	}

	return false
}

func IsErrClient(err error) bool {
	if code := status.Code(err); code == codes.InvalidArgument || code == codes.PermissionDenied || code == codes.AlreadyExists || code == codes.NotFound || code == codes.Aborted {
		return true
	}
	return false
}
