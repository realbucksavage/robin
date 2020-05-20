package traffic

const (
	defaultBindAddr = ":8443"
)

type Config struct {
	BindAddr string `yaml:"bind"`
}
