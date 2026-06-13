package main

import (
	"fmt"
	"go-netdisk/internal/config"
	"go-netdisk/internal/db"
	"go-netdisk/internal/router"
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

	gormDB, err := db.Init(cfg)
	if err != nil {
		log.Fatalf("init db failed: %v", err)
	}

	log.Println("database initialized successfully")

	r := router.New(gormDB, cfg)
	addr := fmt.Sprintf(":%d", cfg.App.Port)

	log.Printf("server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
