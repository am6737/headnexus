package main

import (
	"context"
	"github.com/am6737/headnexus/config"
	"github.com/am6737/headnexus/ports"
	"github.com/am6737/headnexus/service"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	ctx := context.Background()

	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	cfg := config.GenerateConfigTemplate()

	app := service.NewApplication(ctx, &cfg, logger)

	httpSrv := ports.NewHttpHandler(app)

	if err := httpSrv.Start(ctx); err != nil {
		logger.Fatal(err)
	}
}
