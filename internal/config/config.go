package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DSN       string        `env:"DSN" env-required:"true"`
	JWTSecret string        `env:"JWT_SECRET" env-required:"true"`
	JWTTTL    time.Duration `env:"JWT_TTL" env-default:"24h"`
	Addr      string        `env:"ADDR" env-default:":8080"`
}

// LoadConfig reads configuration from the given file, falling back to
// environment variables if the file cannot be read.
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		if envErr := cleanenv.ReadEnv(cfg); envErr != nil {
			return nil, fmt.Errorf("config: %w", envErr)
		}
	}
	return cfg, nil
}
