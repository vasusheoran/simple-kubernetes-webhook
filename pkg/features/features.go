package features

import (
	"os"
	"strconv"
)

var (
	InjectedLabelValue = getStringEnvValueOrDefault("ISTIO_REV_VALUE", "test-rev")
)

func getStringEnvValueOrDefault(name, defaultValue string) string {
	val, _ := os.LookupEnv(name)
	if val == "" {
		return defaultValue
	}

	return val
}

func getBooleanEnvValue(name string, defaultValue bool) bool {
	if val, ok := os.LookupEnv(name); ok {
		booleanVal, err := strconv.ParseBool(val)
		if err != nil {
			return defaultValue
		}
		return booleanVal
	}
	return defaultValue
}
