package system

import (
	"context"
	"eff-gateway/gateway/proxy"

	"log"
	"net/http"
	"os"
	"os/signal"
)

func RunHTTPSv(server *http.Server, cancel context.CancelFunc) {
	err := server.ListenAndServe()
	if err != nil {
		cancel()
		log.Fatalln(err)
	}
}

func InitRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/*", proxy.Forward)
	return mux
}

func OsKill(cancel context.CancelFunc) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s
	cancel()
}
