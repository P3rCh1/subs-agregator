package config

type Config struct {
	Logger *Logger `yaml:"logger"`
	Server *Server `yaml:"server"`
}

type Logger struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
