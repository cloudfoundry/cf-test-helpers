package context_setup

import (
	"encoding/json"
	"fmt"
	"os"
)

//TODO: rename services.config?
type IntegrationConfig struct {
	AppsDomain        string  `json:"apps_domain"`
	ApiEndpoint       string  `json:"api"`
	AdminUser         string  `json:"admin_user"`
	AdminPassword     string  `json:"admin_password"`
	SkipSSLValidation bool    `json:"skip_ssl_validation"`
	TimeoutScale      float64 `json:"timeout_scale"`
}

func LoadConfig(path string, configPtr interface{}) error {
    configFile, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("Loading service integration test config file '%s': %s", path, err.Error())
    }

    decoder := json.NewDecoder(configFile)
    err = decoder.Decode(configPtr)
    if err != nil {
        return fmt.Errorf("Decoding config: %s", err.Error())
    }

    return nil
}

func ValidateConfig(config *IntegrationConfig) error {
	if config.ApiEndpoint == "" {
        return fmt.Errorf("Field 'api' must not be empty")
	}

	if config.AdminUser == "" {
        return fmt.Errorf("Field 'admin_user' must not be empty")
	}

	if config.AdminPassword == "" {
        return fmt.Errorf("Field 'admin_password' must not be empty")
	}

	if config.TimeoutScale == 0 {
		config.TimeoutScale = 1
	} else if config.TimeoutScale < 0 {
        return fmt.Errorf("Field 'timeout_scale' must not be negative (found %d)", config.TimeoutScale)
    }

	return nil
}
