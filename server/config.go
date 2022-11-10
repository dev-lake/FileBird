package main

import (
	"log"

	"gopkg.in/ini.v1"
)

func LoadConfig() *ini.File {
	cfg, err := ini.Load(ConfigPath)
	if err != nil {
		log.Fatalln("Load config.ini failed, err:", err)
	}
	return cfg
}
