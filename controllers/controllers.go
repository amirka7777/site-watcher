package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"site-watcher/models"
	"site-watcher/repository"
	"strconv"
	"time"
)

// структура хранит в себе инструмент для работы с БД
type SiteController struct {
	repositoryHandler repository.SiteRepository
	siteChannel       chan models.Site
}

func NewSiteController(repo repository.SiteRepository, ch chan models.Site) *SiteController {
	return &SiteController{repositoryHandler: repo, siteChannel: ch}
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

	c.siteChannel <- site

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

func (c *SiteController) GetLogsHandler(w http.ResponseWriter, r *http.Request) {

	siteIDStr := r.URL.Query().Get("site_id")
	if siteIDStr == "" {
		http.Error(w, "Параметр site_id обязателен", http.StatusBadRequest)
		return
	}

	siteIDint, err := strconv.Atoi(siteIDStr)
	if err != nil {
		http.Error(w, "поле site_id должно нести в себе число", http.StatusBadRequest)
		return
	}

	var logs []models.CheckLog
	logs, err = c.repositoryHandler.GetLogsBySiteID(siteIDint)
	if err != nil {
		http.Error(w, "Ошибка при получении логов из базы данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)

}
