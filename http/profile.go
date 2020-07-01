package http

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	realworld "github.com/xesina/go-kit-realworld-example-app"
	httpError "github.com/xesina/go-kit-realworld-example-app/http/error"
	"github.com/xesina/go-kit-realworld-example-app/http/middleware"
	"github.com/xesina/go-kit-realworld-example-app/user"
	"net/http"
)

type profileRequest struct {
	username string
	viewerID int64
}

func (req *profileRequest) bind(r *http.Request) error {
	// TODO: handle unexpected errors
	// TODO: move this to a middleware so inject the ID to the context directly
	tk, claims, err := middleware.FromContext(r.Context())
	if err != nil {
		return err
	}

	if tk != nil {
		viewerID := claims["id"].(float64)
		req.viewerID = int64(viewerID)
	}

	username := chi.URLParam(r, "username")
	req.username = username

	if err := req.validate(); err != nil {
		return err
	}

	return nil
}

func (req *profileRequest) validate() error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.username, validation.Required, validation.Length(4, 50)),
		validation.Field(&req.viewerID, is.Int),
	)
}

func (req *profileRequest) endpointRequest() user.ProfileRequest {
	return user.ProfileRequest{
		Username: req.username,
		ViewerID: req.viewerID,
	}
}

type profile struct {
	Username  string          `json:"username"`
	Bio       realworld.Bio   `json:"bio"`
	Image     realworld.Image `json:"image"`
	Following bool            `json:"following"`
}

type profileResponse struct {
	Profile profile `json:"profile"`
}

func newProfileResponse(u *user.ProfileResponse) profileResponse {
	return profileResponse{
		Profile: profile{
			Username:  u.Username,
			Bio:       u.Bio,
			Image:     u.Image,
			Following: u.Following,
		},
	}
}

func (h UserHandler) decodeProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req profileRequest
	if err := req.bind(r); err != nil {
		return nil, err
	}
	er := req.endpointRequest()
	return er, nil
}

func (h UserHandler) encodeProfileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if resp, ok := response.(endpoint.Failer); ok && resp.Failed() != nil {
		httpError.EncodeError(ctx, resp.Failed(), w)
		return nil
	}

	e := response.(user.ProfileResponse)

	hresp := newProfileResponse(&e)

	return jsonResponse(w, hresp, http.StatusOK)
}
