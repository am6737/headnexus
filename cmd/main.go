package main

import (
	"context"
	"flag"
	"github.com/am6737/headnexus/config"
	"github.com/am6737/headnexus/ports"
	"github.com/am6737/headnexus/service"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "配置文件路径")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	if configPath == "" {
		// 从环境变量中获取配置文件路径
		configPath = os.Getenv("CONFIG_PATH")
		if configPath == "" {
			logger.Fatal("未指定配置文件路径，请设置 CONFIG_PATH 环境变量")
		}
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Fatal("无法加载配置文件：", err)
	}

	app := service.NewApplication(ctx, cfg, logger)

	httpSrv := ports.NewHttpHandler(cfg, app)

	if err := httpSrv.Start(ctx); err != nil {
		logger.Fatal(err)
	}
}
