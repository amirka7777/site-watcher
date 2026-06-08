package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"site-watcher/models"
	"site-watcher/repository"
)

// структура хранит в себе инструмент для работы с БД
type SiteController struct {
	repositoryHandler repository.SiteRepository
}

func NewSiteController(repo repository.SiteRepository) *SiteController {
	return &SiteController{repositoryHandler: repo}
}

func (c *SiteController) CreateSiteHandler(w http.ResponseWriter, r *http.Request) {

	var site models.Site

	err := json.NewDecoder(r.Body).Decode(&site)
	if err != nil {
		http.Error(w, "Ошибка при чтении json (отправлены не валидные данные)", http.StatusBadRequest)
		return
	}
	_, err = url.ParseRequestURI(site.URL)
	if err != nil {
		http.Error(w, "Ошибка: неправильно переданный url", http.StatusBadRequest)
		return
	}

	if site.IntervalSeconds <= 0 {
		site.IntervalSeconds = 60
	}

	err = c.repositoryHandler.CreateSite(&site)
	if err != nil {
		http.Error(w, "Ошибка при записи сайта в базу данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)


}