package main

import (
	"context"
	"time"

	libcontext "github.com/otamoe/go-library/context"
	liblogger "github.com/otamoe/go-library/logger"
	libviper "github.com/otamoe/go-library/viper"
	"github.com/otamoe/vptun-server/server"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var logger = liblogger.Get("app")

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.vptun/server")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("vptun-server")
}

func main() {

	var err error
	// viper
	if err = libviper.Parse(); err != nil {
		logger.Panic("vptun-server", zap.Error(err))
	}

	if libviper.PrintDefaults() {
		return
	}

	// 配置 日志
	liblogger.Viper()

	// fx 只显示警告
	// liblogger.SetLevel("fx", zap.WarnLevel)

	// sync
	defer liblogger.Core().Sync()

	ctx := context.Background()
	app := fx.New(
		libcontext.New(ctx),
		liblogger.New(),
		server.New(),
	)

	if app.Err() != nil {
		logger.Panic("App new failed", zap.Error(app.Err()))
	}

	logger.Info("App start...")
	if err = app.Start(context.Background()); err != nil {
		logger.Panic("App start failed", zap.Error(err))
	}
	logger.Info("App started")

	// 停止信号
	<-app.Done()

	// 停止 app
	logger.Info("App stop...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	go func() {
		// 在此接受停止信号
		<-app.Done()
		cancel()
		logger.Panic("App stop failed", zap.Error(err))
	}()
	if err = app.Stop(ctx); err != nil {
		logger.Panic("App stop failed", zap.Error(err))
	}

	logger.Info("App stoped")
	return
}
