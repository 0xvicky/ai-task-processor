package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func EnvInit() {

	//loading .env vars
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Errorf("Fault loading env:%w", envErr)
		return
	}

}
