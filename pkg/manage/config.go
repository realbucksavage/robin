package manage

const (
	defaultBindAddr = ":8089"
	defaultUsername = "guest"
)

type Config struct {
	Bind string               `yaml:"bind"`
	Auth AuthenticationConfig `yaml:"auth"`
}

type AuthenticationConfig struct {
	Username string `yaml:"user"`
	Password string `yaml:"pass"`
}
