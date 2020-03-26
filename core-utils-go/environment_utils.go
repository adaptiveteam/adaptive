package core_utils_go

import (
	"github.com/pkg/errors"
	"os"
)

var (
	// NonEmptyEnv is the environment that is checked for emptiness
	NonEmptyEnv = NonEmptyMap(os.Getenv)
)
// GetNonEmptyEnv reads a single value from non emtpy environment
func GetNonEmptyEnv(name string) string {
	return NonEmptyEnv(name)
}

// NonEmptyMap is a function that converts a mapping function to 
// mapping function that will check value for emptiness
func NonEmptyMap(env func(string)string) func(string)string {
	return func(key string)string {
		value := env(key)
		if value == "" {
			panic(errors.New("Key " + key + " is not defined"))
		}
		return value
	}
}
