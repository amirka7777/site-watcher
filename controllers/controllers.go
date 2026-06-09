package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"site-watcher/models"
	"site-watcher/repository"
	"time"
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
	site.CreatedAt = time.Now()

	err = c.repositoryHandler.CreateSite(&site)
	if err != nil {
		http.Error(w, "Ошибка при записи сайта в базу данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)


}

func (c *SiteController) GetSitesHandler(w http.ResponseWriter, r *http.Request) {

	sites, err := c.repositoryHandler.GetAllSites()
	if err != nil {
		http.Error(w, "Ошибка при взятии информации о сайтах из Базы данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sites)

}