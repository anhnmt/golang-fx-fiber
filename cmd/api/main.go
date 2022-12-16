package main

import (
	"flag"

	"go.uber.org/fx"

	"github.com/xdorro/golang-fx-fiber/internal/server"
	"github.com/xdorro/golang-fx-fiber/pkg/logger"
)

func main() {
	// -log-path is option for command line
	logPath := flag.String("log-path", "logs/data.log", "log file path")
	flag.Parse()

	logger.NewLogger(*logPath)

	var opts []fx.Option
	opts = append(opts, fx.Provide(
		server.NewServer,
	))
	opts = append(opts, fx.Invoke(func(*server.Server) {}))

	app := fx.New(opts...)

	app.Run() // blocks

}
