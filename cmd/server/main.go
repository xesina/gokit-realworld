package main

import (
	"fmt"
	"github.com/xesina/gokit-realworld/article"
	httpTransport "github.com/xesina/gokit-realworld/http"
	"github.com/xesina/gokit-realworld/sqlite"
	"github.com/xesina/gokit-realworld/user"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//in-memory implementation
	//inmemUserRepo := inmem.NewMemUserSaver()
	//inmemArticleRepo := inmem.NewMemArticleRepo()

	s, err := sqlite.NewStorage("./realworld.db")
	if err != nil {
		panic(err)
	}
	s.Migrate()
	userSrv := user.Service{UserRepo: s.NewUserRepository()}
	articleSrv := article.Service{Repo: s.NewArticleRepository()}

	h := httpTransport.MakeHTTPHandler(userSrv, articleSrv)

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
