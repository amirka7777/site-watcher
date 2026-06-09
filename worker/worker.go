package worker

import (
	"log"
	"net/http"
	"site-watcher/models"
	"site-watcher/repository"
	"time"
)

type Worker struct {
	repo repository.SiteRepository
	siteChannel chan models.Site
}

func NewWorker(repo repository.SiteRepository, ch chan models.Site) *Worker {
	return &Worker{repo: repo, siteChannel: ch}
}

func (w *Worker) Start() {

	log.Println("Фоновый воркер успешно запущен")

	sites, err := w.repo.GetAllSites()
	if err != nil {
		log.Printf("Критическая ошибка воркера при старте: не удалось получить сайты: %v", err)
	}

	for _, site := range sites {
		go w.CheckSiteInBackground(site)
	}

	for newSite := range w.siteChannel {
		log.Printf("[Воркер] запускаю мониторинг для %s", newSite.URL)
		go w.CheckSiteInBackground(newSite)
	}

}

func (w *Worker) CheckSiteInBackground(site models.Site) {

	ticker := time.NewTicker(time.Duration(site.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	log.Printf("[Воркер] Запущен мониторинг для %s (интервал: %d сек)", site.URL, site.IntervalSeconds)

	for range ticker.C {
		log.Printf("[Воркер] Проверяю сайт %s...", site.URL)

		res, err := http.Get(site.URL)
		if err != nil {
			log.Printf("❌ Сайт %s лег поспать (Ошибка: %v)", site.URL, err)
			continue
		}

		res.Body.Close()

		if res.StatusCode >= 200 && res.StatusCode < 300 {
			log.Printf("Сайт %s работает (статус: %d)", site.URL, res.StatusCode)
		} else {
			log.Printf("Сайт %s вернул плохой статус-код: %d", site.URL, res.StatusCode)
		}
	}

}