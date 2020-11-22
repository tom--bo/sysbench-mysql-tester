package main

import (
	"errors"
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	Base    BaseConfig
	Backend BackendMySQL
	Target  TargetMySQL
	Scp     ScpConfig
}

type BaseConfig struct {
	SysbenchPath       string `toml:"sysbench_path"`
	SysbenchSenarioDir string `toml:"sysbench_senario_dir"`
}

type BackendMySQL struct {
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	DB       string `toml:"database"`
	Socket   string `toml:"socket"`
}

type TargetMySQL struct {
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	DB       string `toml:"database"`
	Socket   string `toml:"socket"`
}

type ScpConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Path     string `toml:"path"`
}

func readConf() error {
	confFile := homeDir + "/conf/smt.cnf"
	_, err := os.Stat(confFile)
	if os.IsNotExist(err) {
		return errors.New("[Error] Config file doesn't exist at /etc/smt.cnf")
	} else {
		f, err := os.Open(confFile)
		if err != nil {
			return errors.New("[Error] Can't open the config file")
		}
		defer f.Close()

		_, err = toml.DecodeFile(confFile, &conf)
		if err != nil {
			return err
		}
	}
	return nil
}
