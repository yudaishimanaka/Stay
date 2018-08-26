package main

import (
	"log"

	"../models"

	"github.com/BurntSushi/toml"
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	MySQL MySQLConfig `toml:"mysql"`
}

type MySQLConfig struct {
	User 	 string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

func main() {
	var conf Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Fatal(err)
	}

	engine, err := xorm.NewEngine("mysql", conf.MySQL.User + ":" + conf.MySQL.Password+"@/")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := engine.Exec("CREATE DATABASE " + conf.MySQL.Database); err != nil {
		log.Printf("Database already exists.")
	} else {
		engine.Exec("USE " + conf.MySQL.Database)
		engine.CreateTables(models.User{})
		log.Printf("Success initialize.")
	}
}