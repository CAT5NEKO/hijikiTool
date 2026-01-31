package infrastructure

import (
	"errors"
	"os"

	"github.com/CAT5NEKO/hijikiTool/internal/application/ports"
	"github.com/joho/godotenv"
)

type EnvConfigLoader struct {
	envPath string
}

func NewEnvConfigLoader(envPath string) ports.ConfigLoader {
	return &EnvConfigLoader{envPath: envPath}
}

func (l *EnvConfigLoader) Load() (ports.Config, error) {
	if err := godotenv.Load(l.envPath); err != nil {
		return ports.Config{}, err
	}

	config := ports.Config{
		MisskeyHost:  os.Getenv("MISSKEY_HOST"),
		MisskeyToken: os.Getenv("MISSKEY_TOKEN"),
		Visibility:   os.Getenv("MISSKEY_VISIBILITY"),
		LocalOnly:    os.Getenv("MISSKEY_LOCAL_ONLY") == "true",
	}

	if err := l.validate(config); err != nil {
		return ports.Config{}, err
	}

	if config.Visibility == "" {
		config.Visibility = "home"
	}

	return config, nil
}

func (l *EnvConfigLoader) validate(config ports.Config) error {
	if config.MisskeyHost == "" {
		return errors.New("MISSKEY_HOST is required")
	}
	if config.MisskeyToken == "" {
		return errors.New("MISSKEY_TOKEN is required")
	}
	return nil
}
