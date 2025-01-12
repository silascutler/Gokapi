package environment

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

// IsTrue is a placeholder for yes
const IsTrue = "yes"

// IsFalse is a placeholder for no
const IsFalse = "no"

// Environment is a struct containing available env variables
type Environment struct {
	ConfigDir          string
	ConfigFile         string
	ConfigPath         string
	DataDir            string
	AdminName          string
	AdminPassword      string
	WebserverPort      string
	WebserverLocalhost string
	ExternalUrl        string
	RedirectUrl        string
	SaltAdmin          string
	SaltFiles          string
	LengthId           int
	MaxMemory          int
	UseSsl             string
	AwsBucket          string
	AwsRegion          string
	AwsKeyId           string
	AwsKeySecret       string
	AwsEndpoint        string
}

// IsAwsProvided returns true if all required env variables have been set for using AWS S3 / Backblaze
func (e *Environment) IsAwsProvided() bool {
	return e.AwsBucket != "" &&
		e.AwsRegion != "" &&
		e.AwsKeyId != "" &&
		e.AwsKeySecret != ""
}

var defaultValues = defaultsEnvironment{
	CONFIG_DIR:           "config",
	CONFIG_FILE:          "config.json",
	DATA_DIR:             "data",
	LENGTH_ID:            15,
	MAX_MEMORY_UPLOAD_MB: 20,
}

// New parses the env variables
func New() Environment {
	configDir := envString("CONFIG_DIR")
	configFile := envString("CONFIG_FILE")
	configPath := configDir + "/" + configFile

	return Environment{
		ConfigDir:          configDir,
		ConfigFile:         configFile,
		ConfigPath:         configPath,
		DataDir:            envString("DATA_DIR"),
		AdminName:          envString("USERNAME"),
		AdminPassword:      envString("PASSWORD"),
		WebserverPort:      envString("PORT"),
		ExternalUrl:        envString("EXTERNAL_URL"),
		RedirectUrl:        envString("REDIRECT_URL"),
		SaltAdmin:          envString("SALT_ADMIN"),
		SaltFiles:          envString("SALT_FILES"),
		WebserverLocalhost: envBool("LOCALHOST"),
		LengthId:           envInt("LENGTH_ID", 5),
		MaxMemory:          envInt("MAX_MEMORY_UPLOAD_MB", 5),
		UseSsl:             envBool("USE_SSL"),
		AwsBucket:          envString("AWS_BUCKET"),
		AwsRegion:          envString("AWS_REGION"),
		AwsKeyId:           envString("AWS_KEY"),
		AwsKeySecret:       envString("AWS_KEY_SECRET"),
		AwsEndpoint:        envString("AWS_ENDPOINT"),
	}
}

// Looks up an environment variable or returns an empty string
func envString(key string) string {
	val, ok := os.LookupEnv("GOKAPI_" + key)
	if !ok {
		return defaultValues.getString(key)
	}
	return val
}

// Looks up a boolean environment variable, returns either IsFalse or IsTrue
func envBool(key string) string {
	val, ok := os.LookupEnv("GOKAPI_" + key)
	if !ok {
		return ""
	}
	valLower := strings.ToLower(val)
	if valLower == "true" || valLower == "yes" {
		return IsTrue
	}
	if valLower == "false" || valLower == "no" {
		return IsFalse
	}
	return ""
}

// Looks up an environment variable or returns an empty string
func envInt(key string, minValue int) int {
	val, ok := os.LookupEnv("GOKAPI_" + key)
	if !ok {
		return defaultValues.getInt(key)
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return -1
	}
	if intVal < minValue {
		return minValue
	}
	return intVal

}

// Gets the env variable or default value as string
func (structPointer *defaultsEnvironment) getString(name string) string {
	field := reflect.ValueOf(structPointer).Elem().FieldByName(name)
	if field.IsValid() {
		return field.String()
	}
	return ""
}

// Gets the env variable or default value as int
func (structPointer *defaultsEnvironment) getInt(name string) int {
	field := reflect.ValueOf(structPointer).Elem().FieldByName(name)
	if field.IsValid() {
		return int(field.Int())
	}
	return -1
}

type defaultsEnvironment struct {
	CONFIG_DIR           string
	CONFIG_FILE          string
	DATA_DIR             string
	SALT_ADMIN           string
	SALT_FILES           string
	LENGTH_ID            int
	MAX_MEMORY_UPLOAD_MB int
}
