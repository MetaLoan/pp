package utils

import (
	"time"

	"github.com/MetaLoan/pp/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userUUID string) (string, error) {
	claims := jwt.MapClaims{
		"uuid": userUUID,
		"exp":  time.Now().Add(parseDuration(config.GlobalConfig.JWT.Expire)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
}

func parseDuration(d string) time.Duration {
	duration, err := time.ParseDuration(d)
	if err != nil {
		return 24 * time.Hour
	}
	return duration
}
