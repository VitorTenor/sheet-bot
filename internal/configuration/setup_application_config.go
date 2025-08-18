package configuration

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
)

type ApplicationConfig struct {
	Google struct {
		ClientEmail string `yaml:"client_email"`
		PrivateKey  string `yaml:"private_key"`
		ApiUrl      string `yaml:"api_url"`
		SheetId     string `yaml:"sheet_id"`
	} `yaml:"google"`
	WhatsApp struct {
		WebURL     string `yaml:"web_url"`
		GroupName  string `yaml:"group_name"`
		IsArchived bool   `yaml:"is_archived"`
	} `yaml:"whatsapp"`
	Crawler struct {
		UserDataDir string `yaml:"user_data_dir"`
	} `yaml:"crawler"`
	Ai struct {
		IsEnabled bool   `yaml:"is_enabled"`
		ModelURL  string `yaml:"model_url"`
		ModelName string `yaml:"model_name"`
	} `yaml:"ai"`
}

func InitConfig(_ context.Context, path string) (*ApplicationConfig, error) {
	log.Info("loading configuration file")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("configuration file does not exist")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configStr := string(data)
	configStr = replaceEnvVariables(configStr)

	var config ApplicationConfig
	err = yaml.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func replaceEnvVariables(configStr string) string {
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			key := "${" + pair[0] + "}"
			configStr = strings.Replace(configStr, key, pair[1], -1)
		}
	}
	return configStr
}
