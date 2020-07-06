package http

import (
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-kit/kit/log"
	kitTransport "github.com/go-kit/kit/transport"
	transport "github.com/go-kit/kit/transport/http"
	realworld "github.com/xesina/go-kit-realworld-example-app"
	httpError "github.com/xesina/go-kit-realworld-example-app/http/error"
	"github.com/xesina/go-kit-realworld-example-app/http/middleware"
	"net/http"
	"os"
)

func MakeHTTPHandler(userSrv realworld.UserService, articleSrv realworld.ArticleService) http.Handler {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	options := []transport.ServerOption{
		transport.ServerErrorHandler(kitTransport.NewLogErrorHandler(log.With(logger, "component", "HTTP"))),
		transport.ServerErrorEncoder(httpError.EncodeError),
	}

	tokenAuth := middleware.New("HS256", []byte("secret"), nil)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"HEAD", "GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))

	r.Use(chimiddleware.Logger)

	c := Context{
		router:         r,
		jwt:            tokenAuth,
		serverOptions:  options,
		userService:    userSrv,
		articleService: articleSrv,
	}

	RegisterRoutes(c, r)

	return r
}

func wrapHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
