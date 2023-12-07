package main

import (
	"HW-1/controller/stdhttp"
	"HW-1/gates/psg"
	"HW-1/pkg"
	"context"
	"log"
	"os"
	"os/signal"
)

const (
	Addr = "0.0.0.0:9000"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := pkg.MustLoad()

	database, err := psg.NewPsg(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	log.Println("database loaded success")

	ctrl := stdhttp.NewController(Addr, database)

	go func() {
		if err := ctrl.Srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Printf("server started on Addr: %s\n", Addr)

	// Gracefull shutdown

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	if err = ctrl.Srv.Shutdown(ctx); err != nil {
		log.Fatal("wrong shutdown")
	}

	log.Println("success shutdown")
	os.Exit(0)
}
