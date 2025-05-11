package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PgConfig *PgConfig `yaml:"pg_config"`
	SMTP     *SMTP     `yaml:"smtp"`
	Crypto   *Crypto   `yaml:"crypto"`
	Business *Business `yaml:"business"`
}

type PgConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	User        string        `yaml:"user"`
	Password    string        `yaml:"password"`
}

type SMTP struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"post"`
	User  string `yaml:"user"`
	Pass  string `yaml:"pass"`
	Email string `yaml:"email"`
	Name  string `yaml:"name"`
}

type Crypto struct {
	JWTSecret      string `yaml:"jwt_secret"`
	PGPPublicPath  string `yaml:"pgp_public_path"`
	PGPPrivatePath string `yaml:"pgp_private_path"`
	HMACSecret     string `yaml:"hmac_secret"`
}

type Business struct {
	Margin int `yaml:"margin"`
}

func Load(configPath string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(configPath, cfg)
	return cfg, err
}
