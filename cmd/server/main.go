package main

import (
	"fmt"
	httpTransport "github.com/xesina/go-kit-realworld-example-app/http"
	"github.com/xesina/go-kit-realworld-example-app/inmem"
	"github.com/xesina/go-kit-realworld-example-app/user"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	inmemUserRepo := inmem.NewMemUserSaver()
	userSrv := user.Service{Store: inmemUserRepo}

	h := httpTransport.MakeHTTPHandler(userSrv)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errs <- http.ListenAndServe("127.0.0.1:8585", h)
	}()

	<-errs
}
