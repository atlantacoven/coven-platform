package api

import "os"

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

func Env() Environment {
	v := os.Getenv("APP_ENV")
	if v == "" {
		return Development
	}
	return Environment(v)
}

func IsEnv(e Environment) bool {
	return Env() == e
}
