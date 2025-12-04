package config

type Target struct {
	ID     string                 `koanf:"ID"`
	Type   string                 `koanf:"Type"`
	Config map[string]interface{} `koanf:"Config"`
}

type Targets []Target

func (cfg *Config) Targets() (targets Targets, err error) {
	err = cfg.k.Unmarshal("Target", &targets)
	return targets, err
}
