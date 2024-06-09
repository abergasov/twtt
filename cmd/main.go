package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"twtt/internal/config"
	"twtt/internal/logger"
	"twtt/internal/repository/memory"
	"twtt/internal/routes"
	ethnodes "twtt/internal/service/eth_nodes"
	"twtt/internal/service/indexator"
)

var (
	confFile = flag.String("config", "configs/app_conf.yml", "Configs file path")
	appHash  = os.Getenv("GIT_HASH")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	flag.Parse()
	appLog := logger.NewAppSLogger(appHash)

	appLog.Info("app starting", logger.WithString("conf", *confFile))
	appConf, err := config.InitConf(*confFile)
	if err != nil {
		appLog.Fatal("unable to init config", err, logger.WithString("config", *confFile))
	}

	appLog.Info("init repositories")
	repo := memory.NewStorage()

	appLog.Info("init services")
	balancerEL := ethnodes.NewELBalancer(
		appLog,
		appConf.BalancerConf.NormalURLs,
		appConf.BalancerConf.FallbackURLs,
		appConf.BalancerConf.NormalRetries,
		&appConf.BalancerConf.Timeout,
	)

	srv := indexator.NewService(ctx, appLog, balancerEL, repo)
	go srv.Run()

	appLog.Info("init http service")
	appHTTPServer := routes.InitAppRouter(appLog, srv, fmt.Sprintf(":%d", appConf.AppPort))
	defer func() {
		if err = appHTTPServer.Stop(); err != nil {
			appLog.Fatal("unable to stop http service", err)
		}
	}()
	go func() {
		if err = appHTTPServer.Run(); err != nil {
			appLog.Fatal("unable to start http service", err)
		}
	}()

	// register app shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // This blocks the main thread until an interrupt is received
	cancel()
}
