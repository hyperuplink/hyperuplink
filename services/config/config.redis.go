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

func (cfg *Config) RedisAddresses() []string {
	return cfg.k.Strings("Redis.Addresses")
}

func (cfg *Config) RedisMasterName() string {
	return cfg.k.String("Redis.MasterName")
}

func (cfg *Config) RedisUsername() string {
	return cfg.k.String("Redis.Username")
}

func (cfg *Config) RedisPassword() string {
	return cfg.k.String("Redis.Password")
}

func (cfg *Config) RedisDatabase() int {
	return cfg.k.Int("Redis.Database")
}

func (cfg *Config) RedisReset() bool {
	return cfg.k.Bool("Redis.Reset")
}

func (cfg *Config) RedisPoolsize() int {
	return cfg.k.Int("Redis.Poolsize")
}
