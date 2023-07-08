package config

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type (
	SystemConfig struct {
		JwtSecret string `env:"JWT_SECRET" envDefault:"vdFMe2gpco"`

		DBHost   string `env:"DB_HOST" envDefault:"localhost"`
		DBUser   string `env:"DB_USER" envDefault:"postgres"`
		DBPass   string `env:"DB_PASS" envDefault:""`
		DBPort   string `env:"DB_PORT" envDefault:"5432"`
		DBClient string `env:"DB_CLIENT" envDefault:"pgsql"`
		DBName   string `env:"DB_NAME" envDefault:"ppob_api"`

		DigiflazzBaseUrl   string `env:"DIGIFLAZZ_BASE_URL" envDefault:"https://api.digiflazz.com/v1"`
		DigiflazzTopupPath string `env:"DIGIFLAZZ_TOPUP_PATH" envDefault:"/transaction"`
		DigiflazzUsername  string `env:"DIGIFLAZZ_USERNAME" envDefault:"xuvutug1yqvg"`
		DigiflazzApiKey    string `env:"DIGIFLAZZ_API_KEY" envDefault:"dev-00de1590-29b6-11ed-aef6-57b889da7058"`
		DigiflazzTesting   string `env:"DIGIFLAZZ_TESTING" envDefault:"0"`

		XenditBaseUrl      string `env:"XENDIT_BASE_URL" envDefault:"https://api.xendit.co"`
		XenditCreateVAPath string `env:"XENDIT_CREATE_VA_PATH" envDefault:"callback_virtual_accounts"`
		XenditVAExpiration int    `env:"XENDIT_VA_EXPIRE" envDefault:"1"`

		Logger logrus.Logger
	}
)

func LoadEnv() *SystemConfig {
	cfg := SystemConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading system config: %v", err)
	}

	err = env.Parse(&cfg)
	if err != nil {
		log.Printf("Error loading system config: %v", err)
	}

	initLogger(&cfg)

	return &cfg
}

func initLogger(cfg *SystemConfig) {
	cfg.Logger.SetFormatter(&logrus.JSONFormatter{})
	cfg.Logger.SetOutput(os.Stdout)
	cfg.Logger.SetLevel(logrus.InfoLevel)
}
