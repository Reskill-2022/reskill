package config

import (
	"fmt"
	"os"
)

const (
	Port         = "PORT"
	ClientID     = "CLIENT_ID"
	ClientSecret = "CLIENT_SECRET"
	OutlookToken = "OUTLOOK_TOKEN"
)

type Environment map[string]string

func New() (Environment, error) {
	env := make(Environment)
	for _, key := range []string{
		Port,
		ClientID,
		ClientSecret,
		OutlookToken,
	} {
		v, ok := os.LookupEnv(key)
		if !ok {
			return nil, fmt.Errorf("can't find '%s' in environment", key)
		}
		env[key] = v
	}
	return env, nil
}
