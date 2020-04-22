package optredisconfig

type Cfg struct {
	Name         string `json:"name" mapstructure:"name" yaml:"name" toml:"name"`
	Addr         string `json:"addr" mapstructure:"addr" yaml:"addr" toml:"addr"`
	Password     string `json:"password" mapstructure:"password" yaml:"password" toml:"password"`
	DB           int    `json:"db" mapstructure:"db" yaml:"db" toml:"db"`
	MaxRetries   int    `json:"max_retries" mapstructure:"max_retries" yaml:"max_retries" toml:"max_retries"`
	DialTimeout  int    `json:"dial_timeout" mapstructure:"dial_timeout" yaml:"dial_timeout" toml:"dial_timeout"`
	ReadTimeout  int    `json:"read_timeout" mapstructure:"read_timeout" yaml:"read_timeout" toml:"read_timeout"`
	WriteTimeout int    `json:"write_timeout" mapstructure:"write_timeout" yaml:"write_timeout" toml:"write_timeout"`
}

func ByName(list []Cfg, name string) Cfg {
	for _, client := range list {
		if client.Name == name {
			return client
		}
	}
	return Cfg{}
}
