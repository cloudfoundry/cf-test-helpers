package helpers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

type Config struct {
	ApiEndpoint string `json:"api"`
	AppsDomain  string `json:"apps_domain"`

	AdminUser     string `json:"admin_user"`
	AdminPassword string `json:"admin_password"`

	PersistentAppHost      string `json:"persistent_app_host"`
	PersistentAppSpace     string `json:"persistent_app_space"`
	PersistentAppOrg       string `json:"persistent_app_org"`
	PersistentAppQuotaName string `json:"persistent_app_quota_name"`

	SkipSSLValidation bool `json:"skip_ssl_validation"`

	ArtifactsDirectory string `json:"artifacts_directory"`

	DefaultTimeout     time.Duration `json:"default_timeout"`
	CfPushTimeout      time.Duration `json:"cf_push_timeout"`
	LongCurlTimeout    time.Duration `json:"long_curl_timeout"`
	BrokerStartTimeout time.Duration `json:"broker_start_timeout"`

	SyslogDrainPort int    `json:"syslog_drain_port"`
	SyslogIpAddress string `json:"syslog_ip_address"`
}

var loadedConfig *Config

func LoadConfig() Config {
	if loadedConfig == nil {
		loadedConfig = loadConfigJsonFromPath()
	}

	// overwrite config by env
	readEnv("", reflect.ValueOf(loadedConfig))

	if loadedConfig.ApiEndpoint == "" {
		panic("missing configuration 'api'")
	}

	if loadedConfig.AdminUser == "" {
		panic("missing configuration 'admin_user'")
	}

	if loadedConfig.ApiEndpoint == "" {
		panic("missing configuration 'admin_password'")
	}

	return *loadedConfig
}

func loadConfigJsonFromPath() *Config {
	var config *Config = &Config{
		PersistentAppHost:      "CATS-persistent-app",
		PersistentAppSpace:     "CATS-persistent-space",
		PersistentAppOrg:       "CATS-persistent-org",
		PersistentAppQuotaName: "CATS-persistent-quota",

		ArtifactsDirectory: filepath.Join("..", "results"),
	}

	path := configPath()

	configFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(config)
	if err != nil {
		panic(err)
	}

	return config
}

func configPath() string {
	path := os.Getenv("CONFIG")
	if path == "" {
		panic("Must set $CONFIG to point to an integration config .json file.")
	}

	return path
}

func readEnv(prefix string, v reflect.Value) {
	if v.Kind() != reflect.Ptr {
		panic("not ptr")
	}
	value := v.Elem()
	valueType := value.Type()
	fieldsCount := valueType.NumField()
	for i := 0; i < fieldsCount; i++ {
		field := value.Field(i)
		fieldName := valueType.Field(i).Name
		fieldKind := field.Kind()
		env := os.Getenv(prefix + fieldName)
		switch fieldKind {
		case reflect.String:
			if len(env) == 0 {
				continue
			}
			field.SetString(env)
		case reflect.Int, reflect.Int64:
			if len(env) == 0 {
				continue
			}
			val, _ := strconv.ParseInt(env, 0, 64)
			field.SetInt(val)
		case reflect.Bool:
			if len(env) == 0 {
				continue
			}
			val, _ := strconv.ParseBool(env)
			field.SetBool(val)
		case reflect.Map:
			for _, key := range field.MapKeys() {
				new_value := reflect.New(field.MapIndex(key).Type()).Elem()
				new_value.Set(field.MapIndex(key))
				readEnv(prefix+fieldName+"_"+key.String()+"_", new_value.Addr())
				field.SetMapIndex(key, new_value)
			}
		case reflect.Struct:
			readEnv(prefix+fieldName+"_", field.Addr())
		default:
			panic("readProcess undefined. type = " + fieldKind.String())
		}
	}
}
