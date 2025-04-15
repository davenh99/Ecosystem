package auth

import (
	"apps/ecosystem/tools"
	"apps/ecosystem/tools/config"
	"apps/ecosystem/tools/errors"
	"apps/ecosystem/tools/types"
	"context"
	goErrors "errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string
const UserKey contextKey = "userId"
const JWTAccess contextKey = "jwtAccess"

func createJWT(payload jwt.MapClaims, secret []byte, expirationSeconds int64) (string, error) {
	expiration := time.Second * time.Duration(expirationSeconds)

	claims := jwt.MapClaims{
		"expires": time.Now().Add(expiration).Unix(),
	}

	for k, v := range payload {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from the user request
		tokenAccess := r.Header.Get("Authorization")
		// Below is same as above
		// tokenString := getTokenFromRequest(r)

		// validate the JWT
		jwtAccess, err := TokenAccessValidate(tokenAccess)

		// TODO maybe try to clean up all of this code to feel more secure and just make it more logical
		if err == nil {
			if !jwtAccess.Valid {
				log.Printf("invalid token")
				errors.PermissionDenied(w)
				return
			}

			// TODO what the heck is below doing?
			claims := jwtAccess.Claims.(jwt.MapClaims)
			userId := claims["id"].(string)

			u, err := store.GetByID(userId)
			if err != nil {
				log.Printf("failed to get user by id: %v", err)
				errors.PermissionDenied(w)
				return
			}
			
			// set context "userId" to the user Id
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserKey, u.Id)
			r = r.WithContext(ctx)
		} else {
			// TODO should we log that the access token failed validation?
			userId, tokenAccess, tokenRefresh, err := refreshTokens(store, r)
			if err != nil {
				errors.PermissionDenied(w)
				return
			}
			// add new access jwt to the context so we can grab it in the handlerfunc and give to user
			ctx := r.Context()
			ctx = context.WithValue(ctx, JWTAccess, tokenAccess)
			ctx = context.WithValue(ctx, UserKey, userId)
			r = r.WithContext(ctx)
			
			tools.SetCookie(w, tokenRefresh)
		}

		handlerFunc(w, r)
	}
}

func TokenAccessValidate(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Env.JWTAccessSecret), nil
	})
}

func TokenRefreshValidate(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Env.JWTRefreshSecret), nil
	})
}

// TODO can below two funcs be combined?
func GetUserIdFromContext(ctx context.Context) string {
	userId, ok := ctx.Value(UserKey).(string)
	if !ok {
		return ""
	}

	return userId
}

func GetJWTAccessFromContext(ctx context.Context) string {
	jwt, ok := ctx.Value(JWTAccess).(string)
	if !ok {
		return ""
	}

	return jwt
}

func refreshTokens(store types.UserStore, r *http.Request) (id string, tokenAccessStr string, tokenRefreshStr string, err error) {
	// tokenRefresh, err := tools.GetCookie(r)
	cookie, err := r.Cookie("jwt_cookie")
	// log.Printf("did we make it here")
	// TODO this error handling looks nice!? consider implementing it elsewhere? throught the app?
	if err != nil {
		switch {
		case goErrors.Is(err, http.ErrNoCookie):
			log.Printf("refresh cookie not found")
		default:
			log.Println(err)
		}
		return "", "", "", err
	}
	tokenRefresh := cookie.Value
	// validate the JWT
	jwtRefresh, err := TokenRefreshValidate(tokenRefresh)
	if err != nil {
		log.Printf("failed to validate refresh token: %v", err)
		return "", "", "", err
	}

	claims := jwtRefresh.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)

	u, err := store.GetByID(userId)
	if err != nil {
		log.Printf("failed to get user by id: %v", err)
		return
	}

	// TODO maybe more descriptive error messages...
	if !jwtRefresh.Valid {
		log.Printf("invalid refresh token")
		return
	}

	tokenAccess, err := NewUserAccessJWT(u.Id)
	if err != nil {
		log.Printf("error getting new access token: %v", err)
		return "", "", "", err
	}

	tokenRefresh, err = NewUserRefreshJWT(u.Id)
	if err != nil {
		log.Printf("error getting new refresh token: %v", err)
		return "", "", "", err
	}

	return u.Id, tokenAccess, tokenRefresh, nil
}
