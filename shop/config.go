package shop

import (
	"fmt"
	"os"
	"strings"

	"github.com/doutorfinancas/go-mad/core"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Config struct {
	URL        string          `yaml:"url"`
	Build      *ConfigBuild    `yaml:"build,omitempty"`
	AdminApi   *ConfigAdminApi `yaml:"admin_api,omitempty"`
	ConfigDump *ConfigDump     `yaml:"dump,omitempty"`
	Sync       *ConfigSync     `yaml:"sync,omitempty"`
}

type ConfigBuild struct {
	DisableAssetCopy      bool     `yaml:"disable_asset_copy,omitempty"`
	RemoveExtensionAssets bool     `yaml:"remove_extension_assets,omitempty"`
	KeepExtensionSource   bool     `yaml:"keep_extension_source,omitempty"`
	CleanupPaths          []string `yaml:"cleanup_paths,omitempty"`
	Browserslist          string   `yaml:"browserslist,omitempty"`
}

type ConfigAdminApi struct {
	ClientId        string `yaml:"client_id,omitempty"`
	ClientSecret    string `yaml:"client_secret,omitempty"`
	Username        string `yaml:"username,omitempty"`
	Password        string `yaml:"password,omitempty"`
	DisableSSLCheck bool   `yaml:"disable_ssl_check,omitempty"`
}

type ConfigDump struct {
	Rewrite map[string]core.Rewrite `yaml:"rewrite,omitempty"`
	NoData  []string                `yaml:"nodata,omitempty"`
	Ignore  []string                `yaml:"ignore,omitempty"`
	Where   map[string]string       `yaml:"where,omitempty"`
}

type ConfigSync struct {
	Config       []ConfigSyncConfig `yaml:"config"`
	Theme        []ThemeConfig      `yaml:"theme"`
	MailTemplate []MailTemplate     `yaml:"mail_template"`
	Entity       []EntitySync       `yaml:"entity"`
}

type ConfigSyncConfig struct {
	SalesChannel *string                `yaml:"sales_channel,omitempty"`
	Settings     map[string]interface{} `yaml:"settings"`
}

type ThemeConfig struct {
	Name     string                               `yaml:"name"`
	Settings map[string]adminSdk.ThemeConfigValue `yaml:"settings"`
}

type MailTemplate struct {
	Id           string                    `yaml:"id"`
	Translations []MailTemplateTranslation `yaml:"translations"`
}

type EntitySync struct {
	Entity  string                 `yaml:"entity"`
	Exists  *[]interface{}         `yaml:"exists"`
	Payload map[string]interface{} `yaml:"payload"`
}

type MailTemplateTranslation struct {
	Language     string      `yaml:"language"`
	SenderName   string      `yaml:"sender_name"`
	Subject      string      `yaml:"subject"`
	HTML         string      `yaml:"html"`
	Plain        string      `yaml:"plain"`
	CustomFields interface{} `yaml:"custom_fields"`
}

func ReadConfig(fileName string, allowFallback bool) (*Config, error) {
	config := &Config{}

	_, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		if allowFallback {
			return fillEmptyConfig(config), nil
		}

		return nil, fmt.Errorf("cannot find .shopware-project.yml, use shopware-cli project config init to create one")
	}

	if err != nil {
		return nil, err
	}

	fileHandle, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("ReadConfig: %v", err)
	}

	substitutedConfig := os.ExpandEnv(string(fileHandle))
	err = yaml.Unmarshal([]byte(substitutedConfig), &config)

	if err != nil {
		return nil, fmt.Errorf("ReadConfig: %v", err)
	}

	return fillEmptyConfig(config), nil
}

func fillEmptyConfig(c *Config) *Config {
	if c.Build == nil {
		c.Build = &ConfigBuild{}
	}

	return c
}

func NewUuid() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
