package config

type Config struct {
	Logger *Logger `yaml:"logger"`
}

func Default() *Config {
	return &Config{
		Logger: &Logger{
			Level:  "info",
			Format: "text",
		},
	}
}

func ParseFile(path string) (*Config, error) {
	cfg := Default()
	return cfg, nil
}

type Logger struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}
