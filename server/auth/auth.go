package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
)

const (
	tokenExpired         = "Token expired."
	endpointUnauthorized = "Unauthorized access to this endpoint."
	tokenError           = "Bad token."
	claimKey             = "claims"
)

type JWTClaims struct {
	UserID int64  `json:"sub"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

func NewAuthMiddleware(secret string) negroni.Handler {
	middleware := jwtmiddleware.New(jwtmiddleware.Options{
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			w.Write([]byte(tokenError))
		},
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			claims := parseClaims(token.Raw)
			if time.Now().Unix() > claims.Exp {
				return []byte(secret), errors.New(tokenExpired)
			}
			return []byte(secret), nil
		},
		Extractor: jwtmiddleware.FromFirst(jwtmiddleware.FromAuthHeader),
	})

	return negroni.HandlerFunc(middleware.HandlerWithNext)
}

func parseClaims(raw string) JWTClaims {
	claims := JWTClaims{}
	bytes, _ := jwt.DecodeSegment(strings.Split(raw, ".")[1])
	json.Unmarshal(bytes, &claims)
	return claims
}

type AuthorizationMiddleware struct {
	AuthorizationServiceHost string
	Client                   *http.Client
}

func (am AuthorizationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, ok := r.Context().Value("user").(*jwt.Token)
	if !ok {
		Unauthorized(w)
		return
	}
	req, _ := http.NewRequest("GET", am.AuthorizationServiceHost, nil)
	req.Header.Set("Authorization", "Bearer "+token.Raw)
	resp, err := am.Client.Do(req)
	if err != nil {
		log.Println(err)
		Unauthorized(w)
		return
	}
	if resp.StatusCode != 200 {
		Unauthorized(w)
		return
	}
	newRequest := r.WithContext(context.WithValue(r.Context(), claimKey, parseClaims(token.Raw)))
	*r = *newRequest
	next(w, r)
}

func Unauthorized(w http.ResponseWriter) {

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(tokenError))
}
