package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type FileReader struct {
	filename string
	dir      string
}

func NewFileLoader(filename string, dir string) Loader {
	return &FileReader{
		filename: filename,
		dir:      dir,
	}
}

func (r *FileReader) LoadEnv(v viper.Viper) (*viper.Viper, error) {
	filePath := fmt.Sprintf("%s/%s", r.dir, r.filename)

	err := godotenv.Load(filePath)
	if err != nil {
		return nil, err
	}
	v.AutomaticEnv()

	return &v, nil
}
