package auth

import (
	"apps/ecosystem/tools/config"

	"github.com/golang-jwt/jwt/v5"
)

// TODO get user roles and assign different jwt??
func NewUserRefreshJWT(userId string) (string, error) {
	claims := jwt.MapClaims{
		"type": "user",
		"id": userId,
	}
	return createJWT(claims, []byte(config.Env.JWTRefreshSecret), 600) // 10min
}

func NewUserAccessJWT(userId string) (string, error) {
	claims := jwt.MapClaims{
		"type": "user",
		"id": userId,
	}
	return createJWT(claims, []byte(config.Env.JWTAccessSecret), 2592000) // 30 days
}

// func NewAdminJWT(user *types.User) (string, error) {
// 	claims := jwt.MapClaims{
// 		"type": "user",
// 		"id": strconv.Itoa(user.Id),
// 	}
// 	return CreateJWT(claims, []byte(config.Env.JWTSecret))
// }
