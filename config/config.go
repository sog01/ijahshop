package config

import gcfg "gopkg.in/gcfg.v1"

// Config is main entity of package config
type Config struct {
	Storage Storage
}

// Storage is entity of config Storage
type Storage map[string]*struct {
	Host string
}

// New to create instance of config
func New(filePath string) (Config, error) {
	var config Config
	err := gcfg.ReadFileInto(&config, filePath)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
