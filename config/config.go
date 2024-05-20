package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var (
	Global   Config
	ExistENV = false
)

type Config struct {
	DbTablePrefix string `env:"DB_TABLE_PREFIX"`
	DbEngine      string `env:"DB_ENGINE"`
	DbUser        string `env:"DB_USER"`
	DbPassword    string `env:"DB_PASSWORD"`
	DbHost        string `env:"DB_HOST"`
	DbPort        string `env:"DB_PORT"`
	DbName        string `env:"DB_NAME"`
	DbCharset     string `env:"DB_CHARSET"`

	AppIsDebug bool `env:"APP_IS_DEBUG"`

	SessionCookieName string `env:"SESSION_COOKIE_NAME"`
	SessionSecretKey  string `env:"SESSION_SECRET_KEY"`

	SmtpServer   string `env:"SMTP_SERVER"`
	SmtpPort     int    `env:"SMTP_PORT"`
	SmtpUsername string `env:"SMTP_USERNAME"`
	SmtpPassword string `env:"SMTP_PASSWORD"`

	AdminTheme string `env:"ADMIN_THEME"`

	IsResponsive bool `env:"IS_RESPONSIVE"`

	UploadImageResize       string `env:"UPLOAD_IMAGE_RESIZE"`
	UploadImageSizeLimit    int    `env:"UPLOAD_IMAGE_SIZE_LIMIT"`
	UploadImageResizeWidth  int    `env:"UPLOAD_IMAGE_RESIZE_WIDTH"`
	UploadImageResizeHeight int    `env:"UPLOAD_IMAGE_RESIZE_HEIGHT"`
	UploadImageQuality      int    `env:"UPLOAD_IMAGE_QUALITY"`

	CookieDomain string `env:"COOKIE_DOMAIN"`
}

func Load() error {
	setDefaultValue(&Global)

	err := godotenv.Load()
	if err != nil {
		return err
	}

	err = env.Parse(&Global)
	if err != nil {
		return err
	}

	ExistENV = true
	return nil
}

func setDefaultValue(c *Config) {
	c.DbTablePrefix = "g6_"
}
