package configure

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"person-key-saver/internal/app/config"
)

type Config struct {
	WS            config.WSConfig
	Http          config.HttpConfig
	LogLevel      string
	Store         config.DbConfig
	Rmq           config.RmqConfig
	Mqtt          config.MqttConfig
	ClickHouse    config.ClickhouseConfig
	DailyAutoSync bool
}

func NewConfig() *Config {
	return &Config{
		LogLevel: "debug",
	}
}

func (this *Config) SetUp(filename string) error {
	v := viper.New()
	v.SetDefault("LOG_LEVEL", "debug")
	if err := this.configureFromFile(filename, v); err != nil {
		if err := this.configureFromEnvironment(v); err != nil {
			return err
		}
	}

	return nil
}

func (this *Config) configureFromFile(filename string, v *viper.Viper) error {
	v.SetConfigName(filename)
	v.SetConfigType("toml")
	v.AddConfigPath("configs")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&this); err != nil {
		return err
	}
	return nil
}

func (this *Config) ConfigureLogger(logLevel string) (*logrus.Logger, error) {
	logger := logrus.New()

	if level, err := logrus.ParseLevel(logLevel); err != nil {
		return nil, err
	} else {
		logger.SetLevel(level)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
	return logger, nil
}

func (this *Config) configureFromEnvironment(v *viper.Viper) error {
	v.AutomaticEnv()
	var ok bool

	rawLogLevel := v.Get("")
	if this.LogLevel, ok = rawLogLevel.(string); !ok {
		return errors.New("не заполнен уровень логирования")
	}

	{
		rawDbUser := v.Get("*****")
		if this.Store.User, ok = rawDbUser.(string); !ok {
			return errors.New("не заполнен пользователь бд")
		}

		rawDbPwd := v.Get("*****")
		if this.Store.Pwd, ok = rawDbPwd.(string); !ok {
			return errors.New("не заполнен пароль от бд")
		}

		rawDbDsn := v.Get("*****")
		if this.Store.Dsn, ok = rawDbDsn.(string); !ok {
			this.Store.Dsn = ""
		}

		rawDbPort := v.Get("*****")
		if this.Store.Port, ok = rawDbPort.(string); !ok {
			this.Store.Port = ""
		}

		rawDbDatabase := v.Get("*****")
		if this.Store.DataBase, ok = rawDbDatabase.(string); !ok {
			this.Store.DataBase = ""
		}

		rawRmqUser := v.Get("*****")
		if this.Rmq.User, ok = rawRmqUser.(string); !ok {
			return errors.New("не заполнен пользователь для qmqp")
		}

		rawRmqPwd := v.Get("*****")
		if this.Rmq.Pwd, ok = rawRmqPwd.(string); !ok {
			return errors.New("не заполнен пароль для qmqp")
		}

		rawRmqHost := v.Get("*****")
		if this.Rmq.Host, ok = rawRmqHost.(string); !ok {
			return errors.New("не заполнен хост для qmqp")
		}

		rawRmqPort := v.Get("*****")
		if this.Rmq.Port, ok = rawRmqPort.(string); !ok {
			return errors.New("не заполнен порт для qmqp")
		}

		rawRmqVhost := v.Get("*****")
		if this.Rmq.Vhost, ok = rawRmqVhost.(string); !ok {
			this.Rmq.Vhost = ""
		}

		rawChUser := v.Get("*****")
		if this.ClickHouse.User, ok = rawChUser.(string); !ok {
			this.ClickHouse.User = ""
		}

		rawChPwd := v.Get("*****")
		if this.ClickHouse.Pwd, ok = rawChPwd.(string); !ok {
			this.ClickHouse.Pwd = ""
		}

		rawChDsn := v.Get("*****")
		if this.ClickHouse.Dsn, ok = rawChDsn.(string); !ok {
			this.ClickHouse.Dsn = ""
		}

		rawChPort := v.Get("*****")
		if this.ClickHouse.Port, ok = rawChPort.(string); !ok {
			this.ClickHouse.Port = ""
		}

		rawChDatabase := v.Get("*****")
		if this.ClickHouse.Database, ok = rawChDatabase.(string); !ok {
			this.ClickHouse.Database = ""
		}

		rawWSPort := v.Get("*****")
		if this.WS.Port, ok = rawWSPort.(string); !ok {
			this.WS.Port = ""
		}

		rawHttpPort := v.Get("*****")
		if this.Http.Port, ok = rawHttpPort.(string); !ok {
			this.Http.Port = ""
		}

		rawMqttHost := v.Get("*****")
		if this.Mqtt.Host, ok = rawMqttHost.(string); !ok {
			this.Mqtt.Host = ""
		}

		rawMqttClient := v.Get("*****")
		if this.Mqtt.Client, ok = rawMqttClient.(string); !ok {
			this.Mqtt.Client = ""
		}

		rawDailyAutoSync := v.Get("*****")
		if this.DailyAutoSync, ok = rawDailyAutoSync.(bool); !ok {
			this.DailyAutoSync = true
		}
	}

	return nil
}
