package database

type Config struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Database   string `yaml:"database"`
	MaxRetries int    `yaml:"max_retries"`
	LogMode    bool   `yaml:"log_sql"`
}
