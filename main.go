package main

import (
	"log"
	"net/http"
	"site-watcher/controllers"
	"site-watcher/database"
	"site-watcher/repository"
)


func main() {

	db, err := database.InitDB("sites.db")
	if err != nil {
		log.Fatalf("Ошибка при открытии бД: %v", err)
	}

	defer db.Close()

	siteRepo := repository.NewSiteRepository(db)
	siteContr := controllers.NewSiteController(siteRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /sites", siteContr.CreateSiteHandler)
	mux.HandleFunc("GET /sites", siteContr.GetSitesHandler)

	log.Println("Сервер запущен на порту 8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

	

}