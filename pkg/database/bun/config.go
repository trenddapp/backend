package bun

type Config struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Debug    bool   `yaml:"debug"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslMode"`
	Tracing  bool   `yaml:"tracing"`
	User     string `yaml:"user"`
}
