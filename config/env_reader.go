package config

import "github.com/spf13/viper"

type EnvReader struct{}

func NewEnvReader() *EnvReader {
	return &EnvReader{}
}

func (e *EnvReader) LoadEnv(v viper.Viper) (*viper.Viper, error) {
	v.AutomaticEnv()

	return &v, nil
}
