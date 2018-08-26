package main

import (
	"log"

	"./router"

	"github.com/go-xorm/xorm"
	"github.com/BurntSushi/toml"
)

type Config struct {
	MySQL MySQLConfig `toml:"mysql"`
}

type MySQLConfig struct {
	User 	 string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

func init() {
}

func main() {
	var conf Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Fatal(err)
	}

	engine, err := xorm.NewEngine("mysql", conf.MySQL.User + ":" + conf.MySQL.Password+"@/"+conf.MySQL.Database)
	if err != nil {
		log.Fatal(err)
	}

	engine.SetMaxIdleConns(5000)
	engine.SetMaxOpenConns(1000)

	defer engine.Close()

	r := router.Init(engine)
	r.Run(":8888")
}
