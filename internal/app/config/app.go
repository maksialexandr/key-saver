package config

type ApplicationConfig struct {
	Log  LogConfig  `toml:"log"`
	Rmq  RmqConfig  `toml:"rmq"`
	Db   DbConfig   `toml:"db"`
	Mqtt MqttConfig `toml:"mqtt"`
	Http HttpConfig `toml:"http"`
	WS   HttpConfig `toml:"ws"`
}
