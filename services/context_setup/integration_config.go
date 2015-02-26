package context_setup

type IntegrationConfig struct {
	AppsDomain                    string `json:"apps_domain"`
	ApiEndpoint                   string `json:"api"`
	AdminUser                     string `json:"admin_user"`
	AdminPassword                 string `json:"admin_password"`
	CreatePermissiveSecurityGroup bool   `json:"create_permissive_security_group"`
	SkipSSLValidation             bool   `json:"skip_ssl_validation"`
}
