package main

import (
	"log"
	"os"

	"./router"

	"github.com/BurntSushi/toml"
	"github.com/go-xorm/xorm"
)

type Config struct {
	MySQL MySQLConfig `toml:"mysql"`
}

type MySQLConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

func main() {
	// mysql engine create
	var conf Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Fatal(err)
	}

	engine, err := xorm.NewEngine("mysql", conf.MySQL.User+":"+conf.MySQL.Password+"@/"+conf.MySQL.Database)
	if err != nil {
		log.Fatal(err)
	}

	engine.SetMaxIdleConns(5000)
	engine.SetMaxOpenConns(1000)

	defer engine.Close()

	// make directory for user icon
	if _, err := os.Stat("userIcon"); os.IsNotExist(err) {
		err := os.Mkdir("userIcon", 0777)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("successfuly created userIcon dir")
		}
	} else {
		log.Printf("already exists")
	}

	r := router.Init(engine)
	r.Run(":8888")
}
