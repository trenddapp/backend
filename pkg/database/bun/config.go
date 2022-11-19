package bun

type Config struct {
	Debug bool   `yaml:"debug"`
	DSN   string `yaml:"dsn"`
}
