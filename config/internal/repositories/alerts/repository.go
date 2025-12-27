package alerts

import (
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
)

func AddAlert(mailId, resId int64) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec("INSERT INTO Alert (mailId, ressourceId) VALUES (?, ?)", mailId, resId)
	return err
}

func RemoveAlert(mailId, resId int64) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec("DELETE FROM Alert WHERE mailId = ? AND ressourceId = ?", mailId, resId)
	return err
}

// GetAlertsForResource returns all alert configs for a resource
func GetAlertsForResource(resId int64) ([]models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	rows, err := db.Query("SELECT mailId, ressourceId FROM Alert WHERE ressourceId = ?", resId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.MailID, &a.RessourceID); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}
