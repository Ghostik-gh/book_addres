package main

import (
	"HW-1/controller/stdhttp"
	"HW-1/gates/psg"
	"HW-1/pkg"
	"HW-1/pkg/logger"
	"context"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"
)

const (
	Addr = "0.0.0.0:9000"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	logger.SetGlobal(
		zapLogger.With(zap.String("component", "server")),
	)

	cfg := pkg.MustLoad()

	database, err := psg.NewPsg(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	logger.Infof(ctx, "database loaded success")

	ctrl := stdhttp.NewController(ctx, Addr, database)

	go func() {
		if err := ctrl.Srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	logger.Infof(ctx, "server started on Addr: %s\n", Addr)

	// Gracefull shutdown

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	if err = ctrl.Srv.Shutdown(ctx); err != nil {
		logger.Errorf(ctx, "wrong shutdown")
	}

	logger.Infof(ctx, "success shutdown")
	os.Exit(0)
}
