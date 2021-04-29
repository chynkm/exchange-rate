package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/chynkm/ratesdb/datastore"
	"gopkg.in/yaml.v2"
)

const configFile = "config.yml"

var cfg Config

// Config structure from the YAML file
type Config struct {
	DB struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

func init() {
	err := readConfigFile(&cfg)
	if err != nil {
		panic(err.Error())
	}

	datastore.Db, err = sql.Open("mysql", cfg.DB.Username+":"+cfg.DB.Password+"@tcp("+cfg.DB.Host+":"+cfg.DB.Port+")/"+cfg.DB.Name+"?charset=utf8")

	if err != nil {
		panic(err.Error())
	}
}

func main() {
	err := datastore.SaveCurrencyRates()
	if err != nil {
		log.Fatal(err)
	}
}

func readConfigFile(cfg *Config) error {
	f, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}

	return err
}
