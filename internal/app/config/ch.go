package config

type ClickhouseConfig struct {
	User     string
	Pwd      string
	Dsn      string
	Port     string
	Database string
}

func (this *ClickhouseConfig) IsValid() bool {
	return this.User != "" && this.Pwd != "" && this.Dsn != "" && this.Port != "" && this.Database != ""
}
