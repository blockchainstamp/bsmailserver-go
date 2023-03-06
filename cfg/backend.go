package cfg

type BackConfig struct {
	UseMemDB bool `json:"use_mem_db"`
}

func (bc *BackConfig) prepare(cfg, fPath string) error {
	return nil
}
