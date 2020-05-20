package manage

const (
	defaultBindAddr = ":8089"
	defaultUsername = "guest"
	defaultPassword = "guest"
)

type Config struct {
	Bind string `yaml:"bind"`
}

type AuthenticationConfig struct {
	Username string `yaml:"user"`
	Password string `yaml:"password"`
}
