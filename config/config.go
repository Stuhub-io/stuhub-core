package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Loader interface {
	LoadEnv(viper.Viper) (*viper.Viper, error)
}

type Config struct {
	BaseUrl     string
	ServiceName string
	Port        int
	Env         string
	Debug       bool

	RemoteBaseURL string

	AllowedOrigins string
	DBHost         string
	DBPort         string
	DBUser         string
	DBName         string
	DBPass         string
	DBSslMode      bool
	DBDsn          string

	SecretKey                     string
	SendgridKey                   string
	SendgridSetPasswordTemplateId string
	SendgridEmailFrom             string

	HashPwSecretKey string

	// Kafka KafkaConfig
}

type KafkaConfig struct {
	Brokers           string
	Topic             string
	NotificationTopic string
}

func (c *Config) GetCORS() []string {
	cors := strings.Split(c.AllowedOrigins, ";")
	rs := []string{}
	for idx := range cors {
		itm := cors[idx]
		if strings.TrimSpace(itm) != "" {
			rs = append(rs, itm)
		}
	}

	return rs
}

func GetDefaultConfigLoaders() []Loader {
	loaders := []Loader{
		NewEnvReader(),             // Load envs
		NewFileLoader(".env", "."), // Load env from file
	}

	return loaders
}

func generateConfigFromViper(v *viper.Viper) Config {
	return Config{
		BaseUrl:        v.GetString("BASE_URL"),
		Port:           v.GetInt("PORT"),
		Env:            v.GetString("ENV"),
		ServiceName:    v.GetString("SERVICE_NAME"),
		Debug:          v.GetBool("DEBUG"),
		AllowedOrigins: v.GetString("ALLOWED_ORIGINS"),

		RemoteBaseURL: v.GetString("REMOTE_BASE_URL"),

		DBHost:    v.GetString("DB_HOST"),
		DBPort:    v.GetString("DB_PORT"),
		DBUser:    v.GetString("DB_USER"),
		DBName:    v.GetString("DB_NAME"),
		DBPass:    v.GetString("DB_PASS"),
		DBSslMode: v.GetBool("DB_SSL_MODE"),
		DBDsn:     v.GetString("DB_DSN"),

		SecretKey:                     v.GetString("SECRET_KEY"),
		SendgridKey:                   v.GetString("SENDGRID_API_KEY"),
		SendgridSetPasswordTemplateId: v.GetString("SENDGRID_SET_PASSWORD_TEMPLATE_ID"),
		SendgridEmailFrom:             v.GetString("SENDGRID_EMAIL_FROM"),

		HashPwSecretKey: v.GetString("HASH_PW_SECRET_KEY"),
	}
}

func LoadConfig(loaders []Loader) Config {
	v := viper.New()
	v.SetDefault("PORT", "5000")
	v.SetDefault("ENV", "local")
	v.SetDefault("DEBUG", true)

	for idx := range loaders {
		newV, err := loaders[idx].LoadEnv(*v)

		if err == nil {
			v = newV
		}
	}

	return generateConfigFromViper(v)
}

func (c *Config) GetShutdownTimeout() time.Duration {
	return 10 * time.Second
}
