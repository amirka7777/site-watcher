package repository

import (
	"database/sql"
	"fmt"
	"site-watcher/models"
)

type SiteRepository interface {
	CreateSite(site *models.Site) error
	GetAllSites() ([]models.Site, error)
	SaveLog(log *models.CheckLog) error
	GetLogsBySiteID(siteID int) ([]models.CheckLog, error)
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

func (r *SQLiteSiteRepository) GetAllSites() ([]models.Site, error) {

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

func (r *SQLiteSiteRepository) SaveLog(log *models.CheckLog) error {

	query := "INSERT INTO check_logs (site_id, status_code, response_time_ms, is_available) VALUES (?, ?, ?, ?)"
	_, err := r.db.Exec(query, log.SiteID, log.StatusCode, log.ResponseTimeMS, log.IsAvailable)
	if err != nil {
		return fmt.Errorf("Ошибка вставки чек-лога: %v", err)
	}

	return nil

}

func (r *SQLiteSiteRepository) GetLogsBySiteID(siteID int) ([]models.CheckLog, error) {

	query := `SELECT id, site_id, status_code, response_time_ms,
	is_available, checked_at FROM check_logs
	WHERE site_id = ?
	ORDER BY checked_at DESC` // в порядке убывания

	rows, err := r.db.Query(query, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.CheckLog

	for rows.Next() {
		var tmp models.CheckLog

		err := rows.Scan(&tmp.ID, &tmp.SiteID, &tmp.StatusCode, &tmp.ResponseTimeMS, &tmp.IsAvailable, &tmp.CheckedAt)
		if err != nil {
			return nil, err
		}

		logs = append(logs, tmp)
	}
	return logs, nil

}
