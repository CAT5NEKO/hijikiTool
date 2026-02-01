package ports

type Config struct {
	MisskeyHost  string
	MisskeyToken string
	Visibility   string
	LocalOnly    bool
}

type ConfigLoader interface {
	Load() (Config, error)
}
