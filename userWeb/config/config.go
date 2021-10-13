package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	ServerName string        `mapstructure:"name"`
	UserConfig UserSrvConfig `mapstructure:"userSrv"`
}

var UserServerConfig *ServerConfig
