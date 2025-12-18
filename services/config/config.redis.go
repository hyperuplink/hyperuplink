package config

type Redis struct {
	Addrs      []string `koanf:"Addresses"`
	Database   int      `koanf:"Database"`
	MasterName string   `koanf:"MasterName"`
	Username   string   `koanf:"Username"`
	Password   string   `koanf:"Password"`
	Poolsize   int      `koanf:"Poolsize"`
}

func (cfg *Config) Redis() (r Redis, err error) {
	err = cfg.k.Unmarshal("Redis", &r)
	return r, err
}
