package resources

import (
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
)

func GetAll() ([]models.Ressource, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	rows, err := db.Query("SELECT id, name, type FROM Ressource")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []models.Ressource
	for rows.Next() {
		var r models.Ressource
		if err := rows.Scan(&r.ID, &r.Name, &r.Type); err != nil {
			return nil, err
		}
		resources = append(resources, r)
	}
	return resources, nil
}

func GetById(id int64) (*models.Ressource, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	var r models.Ressource
	err = db.QueryRow("SELECT id, name, type FROM Ressource WHERE id = ?", id).Scan(&r.ID, &r.Name, &r.Type)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func Create(name string, resType *string) (int64, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return 0, err
	}
	defer helpers.CloseDB(db)

	res, err := db.Exec("INSERT INTO Ressource (name, type) VALUES (?, ?)", name, resType)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Delete(id int64) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec("DELETE FROM Ressource WHERE id = ?", id)
	return err
}

func Update(id int64, name string, resType *string) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec("UPDATE Ressource SET name = ?, type = ? WHERE id = ?", name, resType, id)
	return err
}
