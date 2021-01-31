package conf

import (
	"log"

	"github.com/imdario/mergo"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// C is the exported global configuration variable
var C Conf

// Conf is the struct containing the configuration of the server
type Conf struct {
	NameServer string `mapstructure:"name_server" form:"name_server"`
	Host       string `mapstructure:"host" form:"host"`
	Port       int    `mapstructure:"port" form:"port"`
	AppendPort bool   `mapstructure:"append_port" form:"append_port"`

	Token string `mapstructure:"token" form:"token"`

	ServeHTTPS bool   `mapstructure:"serve_https" form:"serve_https"`
	SSLCert    string `mapstructure:"ssl_cert" form:"ssl_cert"`
	SSLPrivKey string `mapstructure:"ssl_private_key" form:"ssl_private_key"`

	UploadDir    string  `mapstructure:"upload_dir" form:"upload_dir"`
	DB           string  `mapstructure:"db" form:"db"`
	UniURILength int     `mapstructure:"uniuri_length" form:"uniuri_length"`
	KeyLength    int     `mapstructure:"key_length" form:"key_length"`
	SizeLimit    int64   `mapstructure:"size_limit" form:"size_limit"`
	ViewLimit    int64   `mapstructure:"view_limit" form:"view_limit"`
	DiskQuota    float64 `mapstructure:"disk_quota" form:"disk_quota"`
	LogLevel     string  `mapstructure:"loglevel" form:"loglevel"`

	Stats             bool `mapstructure:"stats" form:"stats"`
	SensitiveMode     bool `mapstructure:"sensitive_mode" form:"sensitive_mode"`
	NoWeb             bool `mapstructure:"no_web" form:"no_web"`
	FullDoc           bool `mapstructure:"fulldoc" form:"fulldoc"`
	AlwaysDownload    bool `mapstructure:"always_download" form:"always_download"`
	DisableEncryption bool `mapstructure:"disable_encryption" form:"disable_encryption"`
	PrometheusEnabled bool `mapstructure:"prometheus_enabled" form:"prometheus_enabled"`
}

// UnparsedConf is the configuration when it's still unparsed properly
type UnparsedConf struct {
	NameServer string `mapstructure:"name_server" form:"name_server"`
	Host       string `mapstructure:"host" form:"host"`
	Port       int    `mapstructure:"port" form:"port"`
	AppendPort bool   `mapstructure:"append_port" form:"append_port"`

	Token string `mapstructure:"token" form:"token"`

	ServeHTTPS bool   `mapstructure:"serve_https" form:"serve_https"`
	SSLCert    string `mapstructure:"ssl_cert" form:"ssl_cert"`
	SSLPrivKey string `mapstructure:"ssl_private_key" form:"ssl_private_key"`

	UploadDir    string  `mapstructure:"upload_dir" form:"upload_dir"`
	DB           string  `mapstructure:"db" form:"db"`
	UniURILength int     `mapstructure:"uniuri_length" form:"uniuri_length"`
	KeyLength    int     `mapstructure:"key_length" form:"key_length"`
	SizeLimit    int64   `mapstructure:"size_limit" form:"size_limit"`
	ViewLimit    int64   `mapstructure:"view_limit" form:"view_limit"`
	DiskQuota    float64 `mapstructure:"disk_quota" form:"disk_quota"`
	LogLevel     string  `mapstructure:"loglevel" form:"loglevel"`

	Stats             bool `mapstructure:"stats" form:"stats"`
	SensitiveMode     bool `mapstructure:"sensitive_mode" form:"sensitive_mode"`
	NoWeb             bool `mapstructure:"no_web" form:"no_web"`
	FullDoc           bool `mapstructure:"fulldoc" form:"fulldoc"`
	AlwaysDownload    bool `mapstructure:"always_download" form:"always_download"`
	DisableEncryption bool `mapstructure:"disable_encryption" form:"disable_encryption"`
	PrometheusEnabled bool `mapstructure:"prometheus_enabled" form:"prometheus_enabled"`
}

// NewDefault returns a Conf instance filled with default values
func NewDefault() Conf {
	return Conf{
		UploadDir:    "up/",
		DB:           "goploader.db",
		Host:         "",
		Port:         8080,
		UniURILength: 10,
		SizeLimit:    20,
		ViewLimit:    5,
		DiskQuota:    0,
		KeyLength:    16,
		LogLevel:     "info",
	}
}

// Validate validates that an unparsed configuration is valid.
func (c *Conf) Validate() map[string]string {
	errors := make(map[string]string)
	if c.NameServer == "" {
		errors["name_server"] = "This field is required."
	}
	if c.ServeHTTPS {
		if c.SSLCert == "" {
			errors["ssl_cert"] = "This field is required if you serve https."
		}
		if c.SSLPrivKey == "" {
			errors["ssl_private_key"] = "This field is required if you serve https."
		}
	}
	return errors
}

// FillDefaults fills the zero value fields in the UnparsedConf with default
// values
func (c *Conf) FillDefaults() error {
	return mergo.Merge(c, NewDefault())
}

// Load loads the given fp (file path) to the C global configuration variable.
func Load(flags *pflag.FlagSet, initial bool) error {
	var err error
	viper.BindPFlags(flags)

	if err = C.FillDefaults(); err != nil {
		return err
	}

	viper.SetEnvPrefix("goploader")
	viper.AutomaticEnv()

	if !initial {
		viper.SetConfigName("conf")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(viper.GetString("conf"))
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				return err
			}
			log.Println(err)
			panic(err)
		}
	}

	err = viper.Unmarshal(&C)

	if err != nil {
		log.Printf("Unable to decode into struct, %v", err)
		return err
	}
	return nil
}
