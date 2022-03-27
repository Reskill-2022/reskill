package config

import (
	"fmt"
	"os"
)

const (
	ProxyCurlApiKey = "PROXYCURL_APIKEY"
	Port            = "PORT"
)

type Environment map[string]string

func New() (Environment, error) {
	env := make(Environment)
	for _, key := range []string{
		Port,
		ProxyCurlApiKey,
	} {
		v, ok := os.LookupEnv(key)
		if !ok {
			return nil, fmt.Errorf("can't find '%s' in environment", key)
		}
		env[key] = v
	}
	return env, nil
}
