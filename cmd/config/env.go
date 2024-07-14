// Package config sets up configuration files and variables
package config

import (
	"net/http"

	apperrors "github.com/Allexsen/Learning-Project/internal/errors"
	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file.
// If the .env file cannot be loaded, it handles the error as a critical failure.
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
