package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	httpError "github.com/xesina/go-kit-realworld-example-app/http/error"
	"net/http"
	"strings"
	"time"
)

var (
	TokenCtxKey = &contextKey{"Token"}
	ErrorCtxKey = &contextKey{"Error"}
)

var (
	ErrUnauthorized = errors.New("jwtauth: token is unauthorized")
	ErrExpired      = errors.New("jwtauth: token is expired")
	ErrNBFInvalid   = errors.New("jwtauth: token nbf validation failed")
	ErrIATInvalid   = errors.New("jwtauth: token iat validation failed")
	ErrNoTokenFound = errors.New("jwtauth: no token found")
	ErrAlgoInvalid  = errors.New("jwtauth: algorithm mismatch")
)

// TODO: get from config
var JWTSecret = []byte("!!SECRET!!")

type JWTAuth struct {
	signKey   interface{}
	verifyKey interface{}
	signer    jwt.SigningMethod
	parser    *jwt.Parser
}

// New creates a JWTAuth authenticator instance that provides middleware handlers
// and encoding/decoding functions for JWT signing.
func New(alg string, signKey interface{}, verifyKey interface{}) *JWTAuth {
	return NewWithParser(alg, &jwt.Parser{}, signKey, verifyKey)
}

// NewWithParser is the same as New, except it supports custom parser settings
// introduced in jwt-go/v2.4.0.
func NewWithParser(alg string, parser *jwt.Parser, signKey interface{}, verifyKey interface{}) *JWTAuth {
	return &JWTAuth{
		signKey:   signKey,
		verifyKey: verifyKey,
		signer:    jwt.GetSigningMethod(alg),
		parser:    parser,
	}
}

func jwtToken(id int64) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, _ := token.SignedString(JWTSecret)
	return t
}

func Verifier(ja *JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Verify(ja, TokenFromHeader)(next)
	}
}

func Verify(ja *JWTAuth, findTokenFns ...func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := VerifyRequest(ja, r, findTokenFns...)
			ctx = NewContext(ctx, token, err)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func VerifyRequest(ja *JWTAuth, r *http.Request, findTokenFns ...func(r *http.Request) string) (*jwt.Token, error) {
	var tokenStr string
	var err error

	// Extract token string from the request by calling token find functions in
	// the order they where provided. Further extraction stops if a function
	// returns a non-empty string.
	for _, fn := range findTokenFns {
		tokenStr = fn(r)
		if tokenStr != "" {
			break
		}
	}

	if tokenStr == "" {
		return nil, ErrNoTokenFound
	}

	// Verify the token
	token, err := ja.Decode(tokenStr)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if verr.Errors&jwt.ValidationErrorExpired > 0 {
				return token, ErrExpired
			} else if verr.Errors&jwt.ValidationErrorIssuedAt > 0 {
				return token, ErrIATInvalid
			} else if verr.Errors&jwt.ValidationErrorNotValidYet > 0 {
				return token, ErrNBFInvalid
			}
		}
		return token, err
	}

	if token == nil || !token.Valid {
		err = ErrUnauthorized
		return token, err
	}

	// Verify signing algorithm
	if token.Method != ja.signer {
		return token, ErrAlgoInvalid
	}

	// Valid!
	return token, nil
}

func (ja *JWTAuth) Encode(claims jwt.Claims) (t *jwt.Token, tokenString string, err error) {
	t = jwt.New(ja.signer)
	t.Claims = claims
	tokenString, err = t.SignedString(ja.signKey)
	t.Raw = tokenString
	return
}

func (ja *JWTAuth) Decode(tokenString string) (t *jwt.Token, err error) {
	t, err = ja.parser.Parse(tokenString, ja.keyFunc)
	if err != nil {
		return nil, err
	}
	return
}

func (ja *JWTAuth) keyFunc(t *jwt.Token) (interface{}, error) {
	if ja.verifyKey != nil {
		return ja.verifyKey, nil
	} else {
		return ja.signKey, nil
	}
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := FromContext(r.Context())

		if err != nil || token == nil || !token.Valid {
			httpError.EncodeError(
				r.Context(),
				httpError.NewError(http.StatusUnauthorized, err),
				w,
			)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

// TokenFromHeader tries to retreive the token string from the
// "Authorization" reqeust header: "Authorization: Token T".
func TokenFromHeader(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if len(token) > 7 && strings.ToUpper(token[0:5]) == "TOKEN" {
		return token[6:]
	}
	return ""
}

func NewContext(ctx context.Context, t *jwt.Token, err error) context.Context {
	ctx = context.WithValue(ctx, TokenCtxKey, t)
	ctx = context.WithValue(ctx, ErrorCtxKey, err)
	return ctx
}

func FromContext(ctx context.Context) (*jwt.Token, jwt.MapClaims, error) {
	token, _ := ctx.Value(TokenCtxKey).(*jwt.Token)

	var claims jwt.MapClaims
	if token != nil {
		if tokenClaims, ok := token.Claims.(jwt.MapClaims); ok {
			claims = tokenClaims
		} else {
			panic(fmt.Sprintf("jwtauth: unknown type of Claims: %T", token.Claims))
		}
	} else {
		claims = jwt.MapClaims{}
	}

	err, _ := ctx.Value(ErrorCtxKey).(error)

	return token, claims, err
}

// UnixTime returns the given time in UTC milliseconds
func UnixTime(tm time.Time) int64 {
	return tm.UTC().Unix()
}

// EpochNow is a helper function that returns the NumericDate time value used by the spec
func EpochNow() int64 {
	return time.Now().UTC().Unix()
}

// ExpireIn is a helper function to return calculated time in the future for "exp" claim
func ExpireIn(tm time.Duration) int64 {
	return EpochNow() + int64(tm.Seconds())
}

// Set issued at ("iat") to specified time in the claims
func SetIssuedAt(claims jwt.MapClaims, tm time.Time) {
	claims["iat"] = tm.UTC().Unix()
}

// Set issued at ("iat") to present time in the claims
func SetIssuedNow(claims jwt.MapClaims) {
	claims["iat"] = EpochNow()
}

// Set expiry ("exp") in the claims
func SetExpiry(claims jwt.MapClaims, tm time.Time) {
	claims["exp"] = tm.UTC().Unix()
}

// Set expiry ("exp") in the claims to some duration from the present time
func SetExpiryIn(claims jwt.MapClaims, tm time.Duration) {
	claims["exp"] = ExpireIn(tm)
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "jwtauth context value " + k.name
}
