package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sayanli/calculator/internal/controller/httpserver"
	"github.com/sayanli/calculator/internal/service"
)

func Run() {
	fmt.Println("Calculator is running...")
	s := service.NewServices()
	httpserver := httpserver.NewHttpServer(s.Calculation)
	go func() {
		log.Print("Http server started on port 8080")
		http.ListenAndServe(":8080", httpserver.Router())
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Print("Server stopped")
}
