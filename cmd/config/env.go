package config

import (
	"net/http"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		apperrors.HandleCriticalError(apperrors.New(
			http.StatusInternalServerError,
			"failed to load .env file",
			apperrors.ErrLoadEnv,
			map[string]interface{}{"details": err.Error()},
		))
	}
}
