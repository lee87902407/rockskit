package rocksdb

func cloneEnvConfig(cfg *EnvConfig) *EnvConfig {
	if cfg == nil {
		return DefaultEnvConfig()
	}
	copy := *cfg
	if copy.LRUSize == "" || copy.BlockSize == "" {
		defaults := DefaultEnvConfig()
		if copy.LRUSize == "" {
			copy.LRUSize = defaults.LRUSize
		}
		if copy.BlockSize == "" {
			copy.BlockSize = defaults.BlockSize
		}
	}
	return &copy
}

func cloneConfig(cfg *Config) *Config {
	if cfg == nil {
		return DefaultConfig()
	}
	copy := *cfg
	defaults := DefaultConfig()
	if copy.WriteBufferSize == "" {
		copy.WriteBufferSize = defaults.WriteBufferSize
	}
	if copy.TargetFileSizeBase == "" {
		copy.TargetFileSizeBase = defaults.TargetFileSizeBase
	}
	if copy.MaxBytesForLevelBase == "" {
		copy.MaxBytesForLevelBase = defaults.MaxBytesForLevelBase
	}
	if copy.Level0FileNumCompactionTrigger == 0 {
		copy.Level0FileNumCompactionTrigger = defaults.Level0FileNumCompactionTrigger
	}
	if copy.CompressionType == "" {
		copy.CompressionType = defaults.CompressionType
	}
	return &copy
}
