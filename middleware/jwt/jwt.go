package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/hanguangbaihuo/sparrow_cloud_go/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type (
	// Token for JWT. Different fields will be used depending on whether you're
	// creating or parsing/verifying a token.
	//
	// A type alias for jwt.Token.
	Token = jwt.Token
	// MapClaims type that uses the map[string]interface{} for JSON decoding
	// This is the default claims type if you don't supply one
	//
	// A type alias for jwt.MapClaims.
	MapClaims = jwt.MapClaims
	// Claims must just have a Valid method that determines
	// if the token is invalid for any supported reason.
	//
	// A type alias for jwt.Claims.
	Claims = jwt.Claims
)

// Shortcuts to create a new Token.
var (
	NewToken           = jwt.New
	NewTokenWithClaims = jwt.NewWithClaims
)

// HS256 and company.
var (
	SigningMethodHS256 = jwt.SigningMethodHS256
	SigningMethodHS384 = jwt.SigningMethodHS384
	SigningMethodHS512 = jwt.SigningMethodHS512
)

// RS256 and company.
var (
	SigningMethodRS256 = jwt.SigningMethodRS256
	SigningMethodRS384 = jwt.SigningMethodRS384
	SigningMethodRS512 = jwt.SigningMethodRS512
)

// ECDSA - EC256 and company.
var (
	SigningMethodES256 = jwt.SigningMethodES256
	SigningMethodES384 = jwt.SigningMethodES384
	SigningMethodES512 = jwt.SigningMethodES512
)

// A function called whenever an error is encountered
type errorHandler func(context.Context, error)

// TokenExtractor is a function that takes a context as input and returns
// either a token or an error.  An error should only be returned if an attempt
// to specify a token was found, but the information was somehow incorrectly
// formed.  In the case where a token is simply not present, this should not
// be treated as an error.  An empty string should be returned in that case.
type TokenExtractor func(context.Context) (string, error)

// Middleware the middleware for JSON Web tokens authentication method
type Middleware struct {
	Config Config
}

// OnError is the default error handler.
// Use it to change the behavior for each error.
// See `Config.ErrorHandler`.
func OnError(ctx context.Context, err error) {
	if err == nil {
		return
	}

	ctx.StopExecution()
	ctx.StatusCode(iris.StatusUnauthorized)
	ctx.JSON(context.Map{"message": err.Error()})
}

// New constructs a new Secure instance with supplied options.
func New(cfg ...Config) *Middleware {

	var c Config
	if len(cfg) == 0 {
		c = Config{}
	} else {
		c = cfg[0]
	}

	// if c.ContextKey == "" {
	// 	c.ContextKey = DefaultContextKey
	// }
	// c.ContextKey = DefaultContextKey

	if c.ErrorHandler == nil {
		c.ErrorHandler = OnError
	}

	if c.Extractor == nil {
		// 默认的 token获取方式
		c.Extractor = FromAuthHeaderToken
	}

	return &Middleware{Config: c}
}

// DefaultJwtMiddleware return default iris jwt middleware
// use like this:
// jwtMiddleware := DefultJwtMiddleware("your_jwt_secret")
// app.Use(jwtMiddleware.Serve)
func DefaultJwtMiddleware(jwtSecret string) *Middleware {
	return New(Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			method, ok := token.Header["alg"].(string)
			if ok {
				return GetSecret(method)
			}
			return []byte(""), fmt.Errorf("[JWT] signing method (alg) is unspecified.")
		},
		CredentialsOptional: true,
	})
}

// Get returns the user (&token) information for this client/request
// func (m *Middleware) Get(ctx context.Context) *jwt.Token {
// 	return ctx.Values().Get(m.Config.ContextKey).(*jwt.Token)
// }

// Serve the middleware's action
func (m *Middleware) Serve(ctx context.Context) {
	_, err := m.CheckJWT(ctx)
	if err != nil {
		m.Config.ErrorHandler(ctx, err)
		return
	}

	// If everything ok then call next.
	ctx.Next()
}

// AutoServe the jwt middleware's action
func AutoServe(ctx context.Context) {
	m := New(Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			method, ok := token.Header["alg"].(string)
			if ok {
				return GetSecret(method)
			}
			return []byte(""), fmt.Errorf("[JWT] signing method (alg) is unspecified.")
		},
		CredentialsOptional: true,
	})
	_, err := m.CheckJWT(ctx)
	if err != nil {
		m.Config.ErrorHandler(ctx, err)
		return
	}

	// If everything ok then call next.
	ctx.Next()
}

// FromAuthHeaderToken is a "TokenExtractor" that takes a give context and extracts
// the JWT token from the Authorization header, header key is "token".
func FromAuthHeaderToken(ctx context.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "token" {
		return "", fmt.Errorf("Authorization header format must be Token {token}")
	}

	return authHeaderParts[1], nil
}

// FromAuthHeader is a "TokenExtractor" that takes a give context and extracts
// the JWT token from the Authorization header.
func FromAuthHeader(ctx context.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

// FromParameter returns a function that extracts the token from the specified
// query string parameter
func FromParameter(param string) TokenExtractor {
	return func(ctx context.Context) (string, error) {
		return ctx.URLParam(param), nil
	}
}

// FromFirst returns a function that runs multiple token extractors and takes the
// first token it finds
func FromFirst(extractors ...TokenExtractor) TokenExtractor {
	return func(ctx context.Context) (string, error) {
		for _, ex := range extractors {
			token, err := ex(ctx)
			if err != nil {
				return "", err
			}
			if token != "" {
				return token, nil
			}
		}
		return "", nil
	}
}

var (
	// ErrTokenMissing is the error value that it's returned when
	// a token is not found based on the token extractor.
	ErrTokenMissing = errors.New("required authorization token not found")

	// ErrTokenInvalid is the error value that it's returned when
	// a token is not valid.
	ErrTokenInvalid = errors.New("token is invalid")

	// ErrTokenExpired is the error value that it's returned when
	// a token value is found and it's valid but it's expired.
	ErrTokenExpired = errors.New("token is expired")
)

var jwtParser = new(jwt.Parser)

// CheckJWT the main functionality, checks for token
func (m *Middleware) CheckJWT(ctx context.Context) (*jwt.Token, error) {
	if !m.Config.EnableAuthOnOptions {
		if ctx.Method() == iris.MethodOptions {
			return nil, nil
		}
	}

	// Use the specified token extractor to extract a token from the request
	token, err := m.Config.Extractor(ctx)

	// If debugging is turned on, log the outcome
	if err != nil {
		utils.LogDebugf(ctx, "Error extracting JWT: %v", err)
		return nil, err
	}

	utils.LogDebugf(ctx, "[JWT] Token extracted: %s", token)

	// If the token is empty...
	if token == "" {
		// Check if it was required
		if m.Config.CredentialsOptional {
			utils.LogDebugf(ctx, "[JWT] No credentials found (CredentialsOptional=true)")
			// No error, just no token (and that is ok given that CredentialsOptional is true)
			return nil, nil
		}

		// If we get here, the required token is missing
		utils.LogDebugf(ctx, "[JWT] Error: No credentials found (CredentialsOptional=false)")
		return nil, ErrTokenMissing
	}
	// Now parse the token
	parsedToken, err := jwtParser.Parse(token, m.Config.ValidationKeyGetter)

	// Check if there was an error in parsing...
	if err != nil {
		utils.LogDebugf(ctx, "[JWT] Error parsing token: %v", err)
		return nil, err
	}

	if m.Config.SigningMethod != nil && m.Config.SigningMethod.Alg() != parsedToken.Header["alg"] {
		err := fmt.Errorf("Expected %s signing method but token specified %s",
			m.Config.SigningMethod.Alg(),
			parsedToken.Header["alg"])
		utils.LogDebugf(ctx, "[JWT] Error validating token algorithm: %v", err)
		return nil, err
	}

	// Check if the parsed token is valid...
	if !parsedToken.Valid {
		utils.LogDebugf(ctx, "[JWT] Token is invalid")
		m.Config.ErrorHandler(ctx, ErrTokenInvalid)
		return nil, ErrTokenInvalid
	}

	utils.LogDebugf(ctx, "[JWT] JWT: %v", parsedToken)

	// only when toke is not empty and valid, we will storage it
	// ctx.Values().Set(RawTokenKey, token)

	// parsedToken.Claims.(MapClaims)
	payload, err := json.Marshal(parsedToken.Claims)
	if err != nil {
		utils.LogDebugf(ctx, "[JWT] json marshal parsed token error: %v", err)
		return parsedToken, err
	}
	b64Payload := base64.StdEncoding.EncodeToString(payload)

	ctx.Request().Header.Set("X-Jwt-Payload", b64Payload)

	// If we get here, everything worked and we can set the
	// user property in context.
	// ctx.Values().Set(m.Config.ContextKey, parsedToken)

	return parsedToken, nil
}
