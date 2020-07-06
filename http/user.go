package http

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	transport "github.com/go-kit/kit/transport/http"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	realworld "github.com/xesina/go-kit-realworld-example-app"
	httpError "github.com/xesina/go-kit-realworld-example-app/http/error"
	"github.com/xesina/go-kit-realworld-example-app/http/middleware"
	"github.com/xesina/go-kit-realworld-example-app/user"
	"io"
	"net/http"
	"time"
)

type UserHandler struct {
	service       realworld.UserService
	jwt           *middleware.JWTAuth
	serverOptions []transport.ServerOption
}

func NewUserHandler(c Context) UserHandler {
	return UserHandler{
		service:       c.userService,
		jwt:           c.jwt,
		serverOptions: c.serverOptions,
	}
}

func (h UserHandler) decodeRegisterRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req userRegisterRequest
	if err := req.bind(r.Body); err != nil {
		return nil, err
	}
	er := req.endpointRequest()
	return er, nil
}

type userRegisterRequest struct {
	User struct {
		Username string `json:"username" valid:"required~username is blank"`
		Email    string `json:"email" valid:"required,email~does not validate as email"`
		Password string `json:"password" valid:"required~First name is blank"`
	} `json:"user"`
}

func (req *userRegisterRequest) bind(r io.Reader) error {
	if e := json.NewDecoder(r).Decode(&req); e != nil {
		return httpError.NewError(http.StatusUnprocessableEntity, httpError.ErrRequestBody)
	}
	if err := req.validate(); err != nil {
		return err
	}
	return nil
}

func (req *userRegisterRequest) validate() error {
	return validation.ValidateStruct(
		&req.User,
		validation.Field(&req.User.Username, validation.Required, validation.Length(4, 50)),
		validation.Field(&req.User.Email, validation.Required, is.Email),
		validation.Field(&req.User.Password, validation.Required, validation.Length(6, 50)),
	)
}

func (req *userRegisterRequest) endpointRequest() user.RegisterRequest {
	return user.RegisterRequest{
		Username: req.User.Username,
		Email:    req.User.Email,
		Password: req.User.Password,
	}
}

type userResponse struct {
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Bio      realworld.Bio   `json:"bio"`
	Image    realworld.Image `json:"image"`
	Token    string          `json:"token"`
}

type userRegisterResponse struct {
	User userResponse `json:"user"`
}

func newUserResponse(u *user.Response) userRegisterResponse {
	return userRegisterResponse{
		User: userResponse{
			Username: u.Username,
			Email:    u.Email,
			Bio:      u.Bio,
			Image:    u.Image,
		},
	}
}

func (h UserHandler) decodeLoginRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req userLoginRequest
	if err := req.bind(r.Body); err != nil {
		return nil, err
	}
	er := req.endpointRequest()
	return er, nil
}

func (h UserHandler) encodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if resp, ok := response.(endpoint.Failer); ok && resp.Failed() != nil {
		httpError.EncodeError(ctx, resp.Failed(), w)
		return nil
	}

	e := response.(user.Response)

	hresp := newUserResponse(&e)

	claims := jwt.MapClaims{
		"id": e.ID,
	}

	middleware.SetIssuedNow(claims)
	middleware.SetExpiryIn(claims, time.Hour*24*5)

	_, tokenString, err := h.jwt.Encode(claims)

	if err != nil {
		return err
	}

	hresp.User.Token = tokenString

	return jsonResponse(w, hresp, http.StatusCreated)
}

type userLoginRequest struct {
	User struct {
		Email    string `json:"email" valid:"required,email~does not validate as email"`
		Password string `json:"password" valid:"required~First name is blank"`
	} `json:"user"`
}

func (req *userLoginRequest) bind(r io.Reader) error {
	if e := json.NewDecoder(r).Decode(&req); e != nil {
		return httpError.NewError(http.StatusUnprocessableEntity, httpError.ErrRequestBody)
	}
	if err := req.validate(); err != nil {
		return err
	}
	return nil
}

func (req *userLoginRequest) validate() error {
	return validation.ValidateStruct(
		&req.User,
		validation.Field(&req.User.Email, validation.Required, is.Email),
		// TODO: check for other way to have validations on password for user update
		//validation.Field(&req.User.Password, validation.Required, validation.Length(6, 50)),
	)
}

func (req *userLoginRequest) endpointRequest() user.LoginRequest {
	return user.LoginRequest{
		Email:    req.User.Email,
		Password: req.User.Password,
	}
}

type userGetRequest struct {
	ID int64
}

func (req *userGetRequest) endpointRequest() user.GetRequest {
	return user.GetRequest{
		ID: req.ID,
	}
}

func (h UserHandler) decodeGetRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	// TODO: handle unexpected errors
	// TODO: move this to a middleware so inject the ID to the context directly
	_, claims, err := middleware.FromContext(r.Context())
	if err != nil {
		return nil, err
	}
	t := claims["id"].(float64)

	id := int64(t)
	req := userGetRequest{
		ID: id,
	}

	er := req.endpointRequest()
	return er, nil
}

func (h UserHandler) decodeUpdateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	// TODO: handle unexpected errors
	// TODO: move this to a middleware so inject the ID to the context directly
	_, claims, err := middleware.FromContext(r.Context())
	if err != nil {
		return nil, err
	}
	t := claims["id"].(float64)
	id := int64(t)

	var req updateRequest
	if err := req.bind(r.Body); err != nil {
		return nil, err
	}
	req.User.ID = id
	er := req.endpointRequest()
	return er, nil
}

type updateRequest struct {
	User struct {
		ID       int64
		Username string          `json:"username"`
		Email    string          `json:"email"`
		Password string          `json:"password"`
		Bio      realworld.Bio   `json:"bio"`
		Image    realworld.Image `json:"image"`
	} `json:"user"`
}

func (req *updateRequest) bind(r io.Reader) error {
	if e := json.NewDecoder(r).Decode(&req); e != nil {
		return httpError.NewError(http.StatusUnprocessableEntity, httpError.ErrRequestBody)
	}

	if err := req.validate(); err != nil {
		return err
	}
	return nil
}

func (req *updateRequest) validate() error {
	return validation.ValidateStruct(
		&req.User,
		validation.Field(&req.User.Username, validation.Required, validation.Length(4, 50)),
		validation.Field(&req.User.Email, validation.Required, is.Email),
	)
}

func (req *updateRequest) endpointRequest() user.UpdateRequest {
	return user.UpdateRequest{
		ID:       req.User.ID,
		Username: req.User.Username,
		Password: req.User.Password,
		Email:    req.User.Email,
		Bio:      req.User.Bio,
		Image:    req.User.Image,
	}
}
