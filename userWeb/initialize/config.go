package initialize

import (
	"github.com/spf13/viper"
	"mxshop-api/userWeb/config"
)

func InitConfig() error {
	v := viper.New()
	v.SetConfigFile("./config.yaml")
	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&config.UserServerConfig)
	if err != nil {
		return err
	}
	return nil
}
