package config

type AuthProvider struct {
	Type   string   `koanf:"Type"`
	Key    string   `koanf:"Key"`
	Secret string   `koanf:"Secret"`
	Scopes []string `koanf:"Scopes"`
}

type AuthProviders []AuthProvider

func (cfg *Config) AuthProviders() (providers AuthProviders, err error) {
	err = cfg.k.Unmarshal("AuthProvider", &providers)
	return providers, err
}
