package middleware

import (
	"fmt"
	"net/http"
	"os"

	gjm "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/grounded042/capacious/utils"
	"github.com/zenazn/goji/web"
)

// JWTMiddleware validates JWTs and does 1 of 3 things based on the JWT:
// 1) If the token is valid, it sets the context variable `UserId` to that of
// the user id from the JWT. The user id is grabbed from the `sub` portion of
// the JWT.
// 2) If the token is not valid, it rejects the request and sends a 400 for bad
// request.
// 3) If no token exists, `UserId` in the context is not set, and the middleware
//  lets the request continue on unhindered.
// It is up to the handlers to act upon the absence or existence of the
// `UserId` variable which represents valid auth.
func JWTMiddleware(c *web.C, h http.Handler) http.Handler {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("GO_JWT_MIDDLEWARE_KEY")), nil
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		// get the token from the header
		authToken, err := gjm.FromAuthHeader(r)

		if err != nil {
			// return bad request if the token is invalid
			w.WriteHeader(http.StatusBadRequest)
			utils.CheckErr(err, "bad request")
		} else if authToken == "" {
			// we still want to process the request, so we are going to serve http,
			// but note that we have not set the `UserId` variable in the context. We
			// still process the request because controllers choose what to do based
			// on the existence or absence of the `UserId` context variable.
			h.ServeHTTP(w, r)
		} else {
			token, err := jwt.Parse(authToken, keyFunc)

			if err != nil {
				if ve, ok := err.(*jwt.ValidationError); ok {
					if ve.Errors&jwt.ValidationErrorMalformed != 0 {
						w.WriteHeader(http.StatusBadRequest)
					} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
						// Token is either expired or not active yet
						w.WriteHeader(http.StatusUnauthorized)
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						utils.CheckErr(err, "something bad happened")
					}
				}
			} else if token.Valid {
				// token is valid, set the user id so other things can use it
				c.Env["UserId"] = token.Claims["sub"]
				h.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}

	return http.HandlerFunc(fn)
}
