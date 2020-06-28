package http

import (
	"context"
	"encoding/json"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	realworld "github.com/xesina/go-kit-realworld-example-app"
	"github.com/xesina/go-kit-realworld-example-app/user"
	"io"
	"net/http"
)

type userRegisterRequest struct {
	User struct {
		Username string `json:"username" valid:"required~username is blank"`
		Email    string `json:"email" valid:"required,email~does not validate as email"`
		Password string `json:"password" valid:"required~First name is blank"`
	} `json:"user"`
}

func (req *userRegisterRequest) bind(r io.Reader) error {
	if e := json.NewDecoder(r).Decode(&req); e != nil {
		return newError(http.StatusUnprocessableEntity, ErrRequestBody)
	}
	if err := req.validate(); err != nil {
		return err
	}
	return nil
}

func (req *userRegisterRequest) validate() error {
	return validation.ValidateStruct(
		&req.User,
		validation.Field(&req.User.Username, validation.Required, validation.Length(5, 50)),
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

func decodeUserRegisterRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req userRegisterRequest
	if err := req.bind(r.Body); err != nil {
		return nil, err
	}
	er := req.endpointRequest()
	return er, nil
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

func newUserRegisterResponse(u *user.Response) userRegisterResponse {
	return userRegisterResponse{
		User: userResponse{
			Username: u.Username,
			Email:    u.Email,
			Bio:      u.Bio,
			Image:    u.Image,
			Token:    "FIX-ME",
		},
	}
}

func encodeUserRegisterResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(user.Response)
	if ok && e.Err != nil {
		encodeError(ctx, e.Err, w)
		return nil
	}
	hresp := newUserRegisterResponse(&e)
	return jsonResponse(w, hresp, http.StatusCreated)
}
