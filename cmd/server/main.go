package main

import (
	"github.com/joho/godotenv" // optional
	"github.com/siddharthray/task-manager-api/internal/config"
	"github.com/siddharthray/task-manager-api/internal/db"
	"github.com/siddharthray/task-manager-api/internal/handler"
	"github.com/siddharthray/task-manager-api/internal/repository"
	"github.com/siddharthray/task-manager-api/internal/router"
	"github.com/siddharthray/task-manager-api/internal/service"
	"log"
)

func main() {
	_ = godotenv.Load() // load .env in dev

	cfg := config.Load()
	sqlDB, err := db.OpenMySQLDB(cfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("error closing DB: %v", err)
		}
	}()

	repo := repository.NewTaskRepository(sqlDB)
	svc := service.NewTaskService(repo)
	hdl := handler.NewTaskHandler(svc)
	r := router.NewRouter(hdl)

	addr := ":" + cfg.HTTPPort
	log.Printf("listening on %s…", addr)
	// handle the error returned by Run()
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
