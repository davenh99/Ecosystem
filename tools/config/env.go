package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port string
	DBUser string
	DBPassword string
	DBAddress string
	DBName string
	JWTRefreshExpirationSeconds int64
	JWTAccessExpirationSeconds int64
	JWTAccessSecret string
	JWTRefreshSecret string
}

var Env Config = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config {
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port: getEnv("PORT", "8080"),
		DBUser: getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBAddress: fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName: getEnv("DB_NAME", "ecosystem"),
		JWTAccessExpirationSeconds: getEnvInt("JWT_ACCESS_EXP", 3600 * 24 * 7),
		JWTRefreshExpirationSeconds: getEnvInt("JWT_REFRESH_EXP", 3600 * 24 * 7),
		JWTAccessSecret: getEnv("JWT_ACCESS_SECRET", "randomsecretthatmightbeabove30characterslong??hrlpm4:)"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "randomsecretthatmightbeabove30characterslong!?rirgopa5&"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			return fallback
		}
		return i
	}

	return fallback
}
