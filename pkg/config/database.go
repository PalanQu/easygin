package config

type DatabaseConfig struct {
	DSN    string `mapstructure:"dsn" defaultvalue:""`
	Driver string `mapstructure:"driver" defaultvalue:"sqlite3"`
}
