package main

import (
	"log"
	"go-netdisk/internal/config"
)

func main() {
	// cfg, err := config.Load("config/dev.yaml")
	cfg, err := config.Load("../../config/dev.yaml")
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	log.Printf("app=%s port=%d mysql=%s:%d/%s",
		cfg.App.Name,
		cfg.App.Port,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DBName,
	)
}
