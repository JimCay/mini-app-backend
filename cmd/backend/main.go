package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
	"tg-backend/config"
	"tg-backend/db"
	"tg-backend/pkg/log"
	"tg-backend/server"
	"time"
)

var app = cli.NewApp()

var (
	Version = "0.0.1"
)

var cliFlag = []cli.Flag{
	config.ConfigFileFlag,
}

func init() {
	app.Action = run
	app.Copyright = "Copyright 2024 Crust Authors"
	app.Name = "Statistic"
	app.Usage = "Statistic"
	app.Authors = []*cli.Author{{Name: "Statistic 2024"}}
	app.Version = Version
	app.EnableBashCompletion = true
	app.Flags = append(app.Flags, cliFlag...)

}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	serv, err := initNode(ctx)
	if err != nil {
		return err
	}
	start(serv)
	return nil
}

func start(serv *server.Server) {
	serv.Start()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	select {
	case <-sigc:
		log.Warn("Interrupt received, shutting down now.")
	}
	err := serv.Shutdown(ctx)
	if err != nil {
		log.Fatal("err shut down ", err)
	}
	os.Exit(0)
}

func initNode(ctx *cli.Context) (*server.Server, error) {
	cfg, err := config.Setup(ctx)
	if err != nil {
		return nil, err
	}

	logWriter, err := log.Setup(cfg.LogConf)
	if err != nil {
		return nil, err
	}

	var storage db.Storage
	storage, err = db.Setup(cfg.DbConf)
	if err != nil {
		return nil, err
	}
	serv := server.NewServer(storage, cfg, logWriter)
	if err != nil {
		return nil, err
	}

	return serv, nil
}
