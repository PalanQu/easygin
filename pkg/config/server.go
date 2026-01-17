package config

type ServerConfig struct {
	Port         string `mapstructure:"port" defaultvalue:"10000"`
	RouterPrefix string `mapstructure:"router_prefix" defaultvalue:"/api"`
	RtartOnError bool   `mapstructure:"restart_on_error" defaultvalue:"false"`
}
