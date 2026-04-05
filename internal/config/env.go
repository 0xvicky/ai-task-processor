package config

import (
	"github.com/joho/godotenv"
)

func EnvInit() error {

	//loading .env vars
	envErr := godotenv.Load()
	if envErr != nil {
		// fmt.Errorf("Fault loading env:%w", envErr)
		return envErr
	}
	return nil
}
