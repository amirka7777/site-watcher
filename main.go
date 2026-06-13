package main

import (
	"log"
	"net/http"
	"site-watcher/controllers"
	"site-watcher/database"
	"site-watcher/models"
	"site-watcher/repository"
	"site-watcher/worker"
)

func main() {

	db, err := database.InitDB("sites.db")
	if err != nil {
		log.Fatalf("Ошибка при открытии бД: %v", err)
	}

	defer db.Close()

	siteRepo := repository.NewSiteRepository(db)
	sharedSiteChannel := make(chan models.Site)
	siteContr := controllers.NewSiteController(siteRepo, sharedSiteChannel)

	siteWorker := worker.NewWorker(siteRepo, sharedSiteChannel)
	go siteWorker.Start()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /sites", siteContr.CreateSiteHandler)
	mux.HandleFunc("GET /sites", siteContr.GetSitesHandler)
	mux.HandleFunc("GET /sites/logs", siteContr.GetLogsHandler)

	log.Println("Сервер запущен на порту 8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

}
