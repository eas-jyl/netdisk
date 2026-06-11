package main

import (
	"go-netdisk/internal/config"
	"go-netdisk/internal/db"
	"log"
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

	if _, err := db.Init(cfg); err != nil {
		log.Fatalf("init db failed: %v", err)
	}

	log.Println("database initialized successfully")
}
