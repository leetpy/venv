package manager

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

type Config struct {
	Debug  bool         `toml:"debug"`
	Server ServerConfig `toml:"server"`
	DB     DB           `toml:"db"`
}

type ServerConfig struct {
	Addr string `toml:"addr"`
	Port int    `toml:"port"`
}

type DB struct {
	DBFile string `toml:"db_file"`
	DBType string `toml:"db_type"`
}

func LoadConfig(cfgFile string) *Config {
	cfg := new(Config)
	var fp *os.File
	var content []byte
	var err error
	if fp, err = os.Open(cfgFile); err != nil {
		panic(err)
	}

	if content, err = ioutil.ReadAll(fp); err != nil {
		panic(err)
	}

	if err = toml.Unmarshal(content, cfg); err != nil {
		panic(err)
	}

	//toml.Decode(cfgFile, &cfg)
	return cfg
}
