package repository

import (
	"database/sql"
	"site-watcher/models"
)


type SiteRepository interface {
	CreateSite(site *models.Site) error
	GetAllSites() ([]models.Site, error)
} 

type SQLiteSiteRepository struct {
	db *sql.DB
}

func NewSiteRepository(db *sql.DB) *SQLiteSiteRepository {
	return &SQLiteSiteRepository{db: db}
}

func (r *SQLiteSiteRepository) CreateSite(site *models.Site) error {

	result, err := r.db.Exec("INSERT INTO sites (url, interval_seconds) VALUES (?, ?)", site.URL, site.IntervalSeconds)
	if err != nil {
		return err
	}

	var id int64
	id, err = result.LastInsertId()
	if err != nil {
		return err
	}
	site.ID = int(id)

	return nil

}

func (r *SQLiteSiteRepository) GetAllSites() ([]models.Site ,error) {

	rows, err := r.db.Query("SELECT id, url, interval_seconds, created_at FROM sites")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var allSites []models.Site
	
	for rows.Next() {
		var tmp models.Site
		err = rows.Scan(&tmp.ID, &tmp.URL, &tmp.IntervalSeconds, &tmp.CreatedAt)
		if err != nil {
			return nil, err
		}
		allSites = append(allSites, tmp)
	}

	return allSites, err

}