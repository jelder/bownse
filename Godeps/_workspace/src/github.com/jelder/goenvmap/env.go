package env

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type EnvMap map[string]string

var (
	ENV = MustLoadEnv()
)

func MustLoadEnv() (env EnvMap) {
	env, err := LoadEnvArrayString(os.Environ())
	if err != nil {
		panic(err)
	}
	return env
}

func MustLoadEnvArrayString(as []string) (env EnvMap) {
	env, err := LoadEnvArrayString(as)
	if err != nil {
		panic(err)
	}
	return env
}

func LoadEnvArrayString(as []string) (env EnvMap, err error) {
	keys := make(map[string]bool)
	env = EnvMap{}
	for _, e := range as {
		pair := strings.SplitN(e, "=", 2)
		if _, dupe := keys[strings.ToUpper(pair[0])]; dupe {
			err = errors.New("Environment contains multiple keys differing only by letter case: " + pair[0])
			return nil, err
		}
		keys[strings.ToUpper(pair[0])] = true
		env[pair[0]] = pair[1]
	}
	return env, err
}

func (env EnvMap) Get(key string, _default string) (val string) {
	val = env[key]
	if val == "" {
		return _default
	}
	return val
}

func (env EnvMap) GetBool(key string) bool {
	val, err := strconv.ParseBool(env[key])
	if err != nil {
		return false
	}
	return val
}

func (env EnvMap) GetNumber(key string, _default float64) float64 {
	val := env[key]
	if val == "" {
		return _default
	}
	float, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return _default
	}
	return float
}

func (env EnvMap) IsSet(key string) bool {
	_, set := env[key]
	return set
}
