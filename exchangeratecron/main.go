package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/chynkm/ratesdb/currencystore"
	"github.com/chynkm/ratesdb/datastore"
	"github.com/chynkm/ratesdb/redisdb"
	"github.com/gomodule/redigo/redis"
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

	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		Database int    `yaml:"database"`
	} `yaml:"redis"`
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

	redisdb.Rdbpool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.Redis.Address)
			if err != nil {
				return nil, err
			}
			if cfg.Redis.Password != "" {
				if _, err := c.Do("AUTH", cfg.Redis.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", cfg.Redis.Database); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}

func main() {
	date, exchangeRates := currencystore.FetchExchangeRates()
	err := datastore.SaveExchangeRates(date, exchangeRates)
	if err != nil {
		log.Fatal(err)
	}

	redisdb.SaveExchangeRates()
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
