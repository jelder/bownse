package env

import (
	"os"
	"strconv"
	"strings"
)

type EnvMap map[string]string

var (
	Env EnvMap
)

func init() {
	Env = LoadEnv()
}

func LoadEnv() (env EnvMap) {
	return LoadEnvArrayString(os.Environ())
}

func LoadEnvArrayString(as []string) EnvMap {
	keys := make(map[string]bool)
	env := EnvMap{}
	for _, e := range as {
		pair := strings.SplitN(e, "=", 2)
		if _, dupe := keys[strings.ToUpper(pair[0])]; dupe {
			panic("Environment contains multiple keys differing only by letter case: " + pair[0])
		}
		env[pair[0]] = pair[1]
	}
	return env
}

func (env EnvMap) GetBool(key string) bool {
	val, err := strconv.ParseBool(env[key])
	if err != nil {
		return false
	}
	return val
}

func (env EnvMap) GetFloat(key string, bitSize int) (float64, error) {
	val := env[key]
	if val == "" {
		return 0, nil
	}
	return strconv.ParseFloat(val, bitSize)
}

func (env EnvMap) GetInt(key string, base int, bitSize int) (int64, error) {
	val := env[key]
	if val == "" {
		return 0, nil
	}
	return strconv.ParseInt(val, base, bitSize)
}

func (env EnvMap) GetUint(key string, base int, bitSize int) (uint64, error) {
	val := env[key]
	if val == "" {
		return 0, nil
	}
	return strconv.ParseUint(val, base, bitSize)
}
