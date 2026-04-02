package rocksdb

func cloneConfig(cfg *Config) *Config {
	if cfg == nil {
		return DefaultConfig()
	}
	copy := *cfg
	return &copy
}
