package app

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	calcgrpc "github.com/sayanli/calculator/internal/controller/grpcserver"
	"github.com/sayanli/calculator/internal/controller/httpserver"
	"github.com/sayanli/calculator/internal/service"
	"google.golang.org/grpc"
)

type App struct {
	log      *slog.Logger
	grpcport int
	httpport int
	golimit  int
}

func NewApp(log *slog.Logger, grpcport int, httpport int, golimit int) *App {
	return &App{
		log:      log,
		grpcport: grpcport,
		httpport: httpport,
		golimit:  golimit,
	}
}

func (a *App) Run() {
	const op = "app.Run"
	log := a.log.With(slog.String("op", op))
	log.Info("Calculator is running...")
	s := service.NewServices(log, a.golimit)
	httpserver := httpserver.NewHttpServer(s.Calculation)

	gRPCServer := grpc.NewServer()
	calcgrpc.RegisterServer(gRPCServer, log, s.Calculation)

	go func() {
		log.Info("Http server started", slog.String("port", fmt.Sprintf("%d", a.httpport)))
		if err := http.ListenAndServe(fmt.Sprintf(":%d", a.httpport), httpserver.Router()); err != nil {
			panic(fmt.Errorf("failed to start http server: %w", err))
		}
	}()

	go func() {
		log.Info("gRPC server started", slog.String("port", fmt.Sprintf("%d", a.grpcport)))
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcport))
		if err != nil {
			panic(fmt.Errorf("failed to listen grpc: %w", err))
		}
		if err := gRPCServer.Serve(lis); err != nil {
			panic(fmt.Errorf("failed to serve grpc: %w", err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("Server stopped")
}
