package http

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	kitTransport "github.com/go-kit/kit/transport"
	transport "github.com/go-kit/kit/transport/http"
	"github.com/labstack/echo/v4"
	"github.com/xesina/go-kit-realworld-example-app/user"
	"net/http"
	"os"
)

func MakeHTTPHandler(s user.Service) http.Handler {
	r := echo.New()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	options := []transport.ServerOption{
		transport.ServerErrorHandler(kitTransport.NewLogErrorHandler(log.With(logger, "component", "HTTP"))),
		transport.ServerErrorEncoder(encodeError),
	}

	r.POST("/users", echo.WrapHandler(transport.NewServer(
		user.RegisterEndpoint(s),
		decodeUserRegisterRequest,
		encodeUserRegisterResponse,
		options...,
	)))

	return r
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(user.Response); ok && e.Err != nil {
		encodeError(ctx, e.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
